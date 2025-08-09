package lucien.security.filesystem

import rego.v1

# Deny dangerous filesystem operations in safe mode
deny contains msg if {
    input.safe_mode == true
    is_dangerous_fs_operation(input.command)
    msg := sprintf("BLOCKED: Dangerous filesystem operation '%s' not allowed in safe mode", [input.command])
}

# Deny operations on sensitive system directories
deny contains msg if {
    input.safe_mode == true
    is_sensitive_directory(input.working_dir)
    not is_read_only_operation(input.command)
    msg := sprintf("BLOCKED: Write operations in sensitive directory '%s' not allowed", [input.working_dir])
}

# Deny operations that could modify system files
deny contains msg if {
    input.safe_mode == true
    some arg in input.args
    is_system_file(arg)
    not is_read_only_operation(input.command)
    msg := sprintf("BLOCKED: Modification of system file '%s' not allowed", [arg])
}

# Deny recursive deletions without explicit confirmation
deny contains msg if {
    input.command == "rm"
    has_recursive_flag(input.args)
    not input.confirmed
    msg := "BLOCKED: Recursive deletion requires confirmation. Use 'rm -rf --confirm' or disable safe mode"
}

# Deny format/partition operations
deny contains msg if {
    is_format_operation(input.command)
    msg := sprintf("BLOCKED: Format/partition operation '%s' not allowed in safe mode", [input.command])
}

# Helper functions

is_dangerous_fs_operation(cmd) if {
    cmd in [
        "rm", "rmdir", "del", "rd",
        "format", "fdisk", "mkfs",
        "dd", "shred", "wipe",
        "chmod", "chown", "chgrp",
        "mv", "move", "ren", "rename"
    ]
}

is_sensitive_directory(path) if {
    sensitive_dirs := [
        "/", "/bin", "/sbin", "/usr", "/etc", "/boot", "/sys", "/proc",
        "C:\\Windows", "C:\\Program Files", "C:\\Program Files (x86)",
        "/System", "/Library", "/Applications"
    ]
    some dir in sensitive_dirs
    startswith(path, dir)
}

is_system_file(path) if {
    system_files := [
        "/etc/passwd", "/etc/shadow", "/etc/hosts",
        "/boot/vmlinuz", "/boot/initrd",
        "C:\\Windows\\System32", "C:\\boot.ini",
        "/System/Library", "/Library/LaunchDaemons"
    ]
    some file in system_files
    startswith(path, file)
}

is_read_only_operation(cmd) if {
    cmd in [
        "ls", "dir", "cat", "type", "head", "tail",
        "grep", "find", "locate", "which", "whereis",
        "file", "stat", "du", "df", "mount"
    ]
}

has_recursive_flag(args) if {
    some arg in args
    arg in ["-r", "-R", "--recursive", "-rf", "-Rf"]
}

is_format_operation(cmd) if {
    cmd in [
        "format", "mkfs", "mkfs.ext4", "mkfs.ntfs",
        "fdisk", "parted", "gparted", "diskpart"
    ]
}

# Allow operations with specific whitelist patterns
allow if {
    input.safe_mode == true
    is_whitelisted_operation(input)
}

is_whitelisted_operation(op) if {
    # Allow operations in user home directory
    startswith(op.working_dir, "/home/")
    not contains(op.working_dir, "../")
}

is_whitelisted_operation(op) if {
    # Allow operations in Windows user directory
    regex.match(`^C:\\Users\\[^\\]+`, op.working_dir)
    not contains(op.working_dir, "..")
}

is_whitelisted_operation(op) if {
    # Allow operations in temporary directories
    op.working_dir in ["/tmp", "/var/tmp", "C:\\Temp", "C:\\Windows\\Temp"]
}

# Rate limiting for filesystem operations
deny contains msg if {
    input.safe_mode == true
    is_bulk_operation(input)
    not input.rate_limited
    msg := "BLOCKED: Bulk filesystem operations require rate limiting in safe mode"
}

is_bulk_operation(op) if {
    op.command in ["find", "grep", "rm"]
    count(op.args) > 10
}

is_bulk_operation(op) if {
    op.command == "rm"
    some arg in op.args
    contains(arg, "*")
}

# Audit logging requirement
warn contains msg if {
    input.safe_mode == true
    is_dangerous_fs_operation(input.command)
    not input.audit_logged
    msg := sprintf("WARNING: Dangerous operation '%s' should be audit logged", [input.command])
}

# Path traversal protection
deny contains msg if {
    some arg in input.args
    contains(arg, "../")
    msg := sprintf("BLOCKED: Path traversal detected in argument '%s'", [arg])
}

deny contains msg if {
    some arg in input.args
    contains(arg, "..\\")
    msg := sprintf("BLOCKED: Path traversal detected in argument '%s'", [arg])
}

# Symlink protection
deny contains msg if {
    input.safe_mode == true
    input.command in ["ln", "mklink"]
    some arg in input.args
    is_sensitive_directory(arg)
    msg := sprintf("BLOCKED: Creating symlinks to sensitive directory '%s' not allowed", [arg])
}