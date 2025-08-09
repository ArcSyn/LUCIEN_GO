package lucien.security.network

import rego.v1

# Deny outbound connections to sensitive ports
deny contains msg if {
    input.safe_mode == true
    is_network_command(input.command)
    some arg in input.args
    port := extract_port(arg)
    is_sensitive_port(port)
    msg := sprintf("BLOCKED: Connection to sensitive port %d not allowed in safe mode", [port])
}

# Deny connections to localhost privileged ports
deny contains msg if {
    input.safe_mode == true
    is_network_command(input.command)
    some arg in input.args
    is_localhost_connection(arg)
    port := extract_port(arg)
    port < 1024
    msg := sprintf("BLOCKED: Connection to localhost privileged port %d not allowed", [port])
}

# Deny connections to private network ranges in safe mode
deny contains msg if {
    input.safe_mode == true
    is_network_command(input.command)
    some arg in input.args
    host := extract_host(arg)
    is_private_network(host)
    msg := sprintf("BLOCKED: Connection to private network '%s' not allowed in safe mode", [host])
}

# Deny dangerous network scanning operations
deny contains msg if {
    input.safe_mode == true
    is_network_scan_command(input.command)
    msg := sprintf("BLOCKED: Network scanning with '%s' not allowed in safe mode", [input.command])
}

# Deny server/listener operations on privileged ports
deny contains msg if {
    input.safe_mode == true
    is_server_command(input.command)
    some arg in input.args
    port := extract_port(arg)
    port < 1024
    msg := sprintf("BLOCKED: Server operations on privileged port %d not allowed", [port])
}

# Deny download operations without verification
deny contains msg if {
    input.safe_mode == true
    is_download_command(input.command)
    not has_verification_flag(input.args)
    msg := sprintf("BLOCKED: Download with '%s' requires verification in safe mode", [input.command])
}

# Deny raw socket operations
deny contains msg if {
    input.safe_mode == true
    is_raw_socket_command(input.command)
    msg := sprintf("BLOCKED: Raw socket operations with '%s' not allowed in safe mode", [input.command])
}

# Deny tunnel/proxy operations
deny contains msg if {
    input.safe_mode == true
    is_tunnel_command(input.command)
    msg := sprintf("BLOCKED: Tunnel/proxy operations with '%s' not allowed in safe mode", [input.command])
}

# Helper functions

is_network_command(cmd) if {
    cmd in [
        "curl", "wget", "nc", "ncat", "netcat",
        "telnet", "ssh", "scp", "rsync", "ftp",
        "ping", "traceroute", "tracert", "nslookup", "dig"
    ]
}

is_network_scan_command(cmd) if {
    cmd in [
        "nmap", "masscan", "zmap", "unicornscan",
        "hping", "hping3", "arping", "fping"
    ]
}

is_server_command(cmd) if {
    cmd in [
        "nc", "ncat", "netcat", "socat",
        "python", "python3", "node", "ruby",
        "php", "perl"
    ]
}

is_download_command(cmd) if {
    cmd in ["curl", "wget", "fetch", "aria2c"]
}

is_raw_socket_command(cmd) if {
    cmd in ["hping", "hping3", "scapy", "nemesis"]
}

is_tunnel_command(cmd) if {
    cmd in [
        "ssh", "autossh", "stunnel", "socat",
        "ngrok", "localtunnel", "chisel"
    ]
}

extract_port(arg) := port if {
    # Extract port from host:port format
    contains(arg, ":")
    parts := split(arg, ":")
    port := to_number(parts[count(parts) - 1])
}

extract_port(arg) := port if {
    # Extract port from -p flag
    startswith(arg, "-p")
    port_str := substring(arg, 2, -1)
    port := to_number(port_str)
}

extract_port(arg) := 80 if {
    # Default HTTP port
    contains(arg, "http://")
    not contains(arg, ":")
}

extract_port(arg) := 443 if {
    # Default HTTPS port
    contains(arg, "https://")
    not contains(arg, ":")
}

extract_host(arg) := host if {
    contains(arg, "://")
    url_parts := split(arg, "://")
    host_port := split(url_parts[1], "/")[0]
    host := split(host_port, ":")[0]
}

extract_host(arg) := host if {
    contains(arg, ":")
    not contains(arg, "://")
    host := split(arg, ":")[0]
}

extract_host(arg) := arg if {
    not contains(arg, ":")
    not contains(arg, "://")
}

is_sensitive_port(port) if {
    sensitive_ports := [
        22, 23, 25, 53, 80, 110, 143, 443, 993, 995,
        135, 137, 138, 139, 445, 1433, 1521, 3306, 3389, 5432,
        6379, 11211, 27017, 50070
    ]
    port in sensitive_ports
}

is_localhost_connection(arg) if {
    localhost_patterns := ["127.0.0.1", "::1", "localhost"]
    some pattern in localhost_patterns
    contains(arg, pattern)
}

is_private_network(host) if {
    # RFC 1918 private networks
    private_ranges := [
        "10.", "172.16.", "172.17.", "172.18.", "172.19.",
        "172.20.", "172.21.", "172.22.", "172.23.", "172.24.",
        "172.25.", "172.26.", "172.27.", "172.28.", "172.29.",
        "172.30.", "172.31.", "192.168."
    ]
    some range in private_ranges
    startswith(host, range)
}

is_private_network(host) if {
    # Link-local addresses
    startswith(host, "169.254.")
}

has_verification_flag(args) if {
    verification_flags := [
        "--check-certificate", "--verify", "--checksum",
        "--hash", "--signature", "-k", "--insecure"
    ]
    some flag in verification_flags
    flag in args
}

# Allow specific whitelisted network operations
allow if {
    input.safe_mode == true
    is_whitelisted_network_operation(input)
}

is_whitelisted_network_operation(op) if {
    # Allow ping to public DNS servers
    op.command == "ping"
    some arg in op.args
    arg in ["8.8.8.8", "1.1.1.1", "9.9.9.9"]
}

is_whitelisted_network_operation(op) if {
    # Allow HTTPS to trusted domains
    op.command in ["curl", "wget"]
    some arg in op.args
    contains(arg, "https://")
    is_trusted_domain(extract_host(arg))
}

is_trusted_domain(host) if {
    trusted_domains := [
        "github.com", "gitlab.com", "bitbucket.org",
        "pypi.org", "npmjs.com", "crates.io",
        "archive.ubuntu.com", "security.ubuntu.com",
        "download.microsoft.com", "packages.microsoft.com"
    ]
    host in trusted_domains
}

is_trusted_domain(host) if {
    # Allow subdomains of trusted domains
    trusted_suffixes := [
        ".github.com", ".githubusercontent.com",
        ".ubuntu.com", ".microsoft.com", ".google.com"
    ]
    some suffix in trusted_suffixes
    endswith(host, suffix)
}

# Rate limiting for network operations
deny contains msg if {
    input.safe_mode == true
    is_bulk_network_operation(input)
    not input.rate_limited
    msg := "BLOCKED: Bulk network operations require rate limiting in safe mode"
}

is_bulk_network_operation(op) if {
    op.command in ["curl", "wget"]
    count(op.args) > 5
}

is_bulk_network_operation(op) if {
    op.command in ["ping", "traceroute"]
    some arg in op.args
    arg in ["-c", "--count"]
    # Check if count is high (simplified check)
}

# Audit logging for network operations
warn contains msg if {
    input.safe_mode == true
    is_network_command(input.command)
    not input.audit_logged
    msg := sprintf("WARNING: Network operation '%s' should be audit logged", [input.command])
}

# DNS security checks
deny contains msg if {
    input.safe_mode == true
    input.command in ["nslookup", "dig", "host"]
    some arg in input.args
    is_suspicious_domain(arg)
    msg := sprintf("BLOCKED: DNS lookup for suspicious domain '%s' not allowed", [arg])
}

is_suspicious_domain(domain) if {
    suspicious_patterns := [
        ".tk", ".ml", ".ga", ".cf",  # Free TLDs often used maliciously
        "tempmail", "guerrilla", "10minutemail",  # Temporary email patterns
        "bit.ly", "tinyurl", "t.co"  # URL shorteners
    ]
    some pattern in suspicious_patterns
    contains(domain, pattern)
}

# Protocol security
deny contains msg if {
    input.safe_mode == true
    is_insecure_protocol(input.command, input.args)
    msg := sprintf("BLOCKED: Insecure protocol usage with '%s' not allowed in safe mode", [input.command])
}

is_insecure_protocol(cmd, args) if {
    cmd in ["curl", "wget"]
    some arg in args
    contains(arg, "http://")  # HTTP instead of HTTPS
    not contains(arg, "localhost")  # Allow HTTP to localhost
}

is_insecure_protocol(cmd, args) if {
    cmd == "ftp"  # Plain FTP is insecure
}

is_insecure_protocol(cmd, args) if {
    cmd == "telnet"  # Telnet is insecure
}