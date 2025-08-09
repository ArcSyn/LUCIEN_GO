"""
DesignerAgent - AI-powered UI code generation
"""

import re
import os
import json
from typing import Dict, List, Optional
from pathlib import Path

class DesignerAgent:
    """AI-powered design agent that generates UI code from natural language"""
    
    def __init__(self):
        self.output_dir = Path("snapmethod/exports")
        self.component_templates = {
            'react': self._get_react_templates(),
            'vue': self._get_vue_templates(),
            'angular': self._get_angular_templates()
        }
        self.css_frameworks = {
            'tailwind': self._get_tailwind_classes(),
            'bootstrap': self._get_bootstrap_classes(),
            'material': self._get_material_classes()
        }
        
    def generate(self, prompt: str) -> str:
        """Generate UI code from natural language prompt"""
        
        # Parse the prompt to extract design requirements
        requirements = self._parse_prompt(prompt)
        
        # Determine the best framework and styling approach
        framework = self._determine_framework(requirements)
        styling = self._determine_styling(requirements)
        
        # Generate the code
        code = self._generate_component(requirements, framework, styling)
        
        # Save to file
        file_path = self._save_code(code, requirements['component_name'], framework)
        
        return code
    
    def _parse_prompt(self, prompt: str) -> Dict:
        """Parse the design prompt to extract requirements"""
        prompt_lower = prompt.lower()
        
        requirements = {
            'component_name': self._extract_component_name(prompt),
            'component_type': self._extract_component_type(prompt_lower),
            'colors': self._extract_colors(prompt_lower),
            'layout': self._extract_layout(prompt_lower),
            'features': self._extract_features(prompt_lower),
            'styling': self._extract_styling_preferences(prompt_lower)
        }
        
        return requirements
    
    def _extract_component_name(self, prompt: str) -> str:
        """Extract component name from prompt"""
        # Look for explicit component names
        name_patterns = [
            r'(\w+)\s+(?:page|component|form|modal|card)',
            r'(?:create|build|make)\s+a?\s*(\w+)',
            r'(\w+)\s+(?:interface|ui|screen)'
        ]
        
        for pattern in name_patterns:
            match = re.search(pattern, prompt.lower())
            if match:
                name = match.group(1).capitalize()
                if name not in ['a', 'the', 'an', 'with', 'for']:
                    return f"{name}Component"
        
        # Default fallback
        if 'login' in prompt.lower():
            return "LoginComponent"
        elif 'form' in prompt.lower():
            return "FormComponent"
        elif 'card' in prompt.lower():
            return "CardComponent"
        elif 'modal' in prompt.lower():
            return "ModalComponent"
        elif 'nav' in prompt.lower():
            return "NavigationComponent"
        else:
            return "CustomComponent"
    
    def _extract_component_type(self, prompt: str) -> str:
        """Determine the type of component to generate"""
        if any(word in prompt for word in ['login', 'signin', 'auth']):
            return 'login_form'
        elif any(word in prompt for word in ['signup', 'register']):
            return 'signup_form'
        elif 'form' in prompt:
            return 'form'
        elif any(word in prompt for word in ['card', 'profile']):
            return 'card'
        elif any(word in prompt for word in ['nav', 'navigation', 'menu']):
            return 'navigation'
        elif any(word in prompt for word in ['modal', 'dialog', 'popup']):
            return 'modal'
        elif any(word in prompt for word in ['button', 'btn']):
            return 'button'
        elif any(word in prompt for word in ['table', 'list']):
            return 'table'
        else:
            return 'generic'
    
    def _extract_colors(self, prompt: str) -> Dict[str, str]:
        """Extract color preferences from prompt"""
        colors = {
            'primary': '#3b82f6',  # blue
            'secondary': '#64748b', # gray
            'accent': '#10b981',   # green
            'background': '#ffffff',
            'text': '#1f2937'
        }
        
        # Dark theme detection
        if any(word in prompt for word in ['dark', 'black', 'night']):
            colors.update({
                'background': '#1f2937',
                'text': '#f9fafb',
                'secondary': '#374151'
            })
        
        # Color-specific detection
        color_map = {
            'blue': '#3b82f6',
            'red': '#ef4444',
            'green': '#10b981', 
            'purple': '#8b5cf6',
            'yellow': '#f59e0b',
            'pink': '#ec4899',
            'indigo': '#6366f1',
            'orange': '#f97316'
        }
        
        for color_name, hex_value in color_map.items():
            if color_name in prompt:
                colors['primary'] = hex_value
        
        # Neon/glow effect detection
        if any(word in prompt for word in ['neon', 'glow', 'bright']):
            colors['accent'] = '#00ffff'  # cyan glow
            
        return colors
    
    def _extract_layout(self, prompt: str) -> str:
        """Determine layout type"""
        if any(word in prompt for word in ['center', 'centered']):
            return 'centered'
        elif any(word in prompt for word in ['sidebar', 'side']):
            return 'sidebar'
        elif any(word in prompt for word in ['grid', 'columns']):
            return 'grid'
        elif any(word in prompt for word in ['flex', 'horizontal']):
            return 'flex'
        else:
            return 'default'
    
    def _extract_features(self, prompt: str) -> List[str]:
        """Extract specific features mentioned in prompt"""
        features = []
        
        feature_keywords = {
            'validation': ['validation', 'validate', 'error'],
            'animation': ['animation', 'animate', 'transition'],
            'responsive': ['responsive', 'mobile', 'tablet'],
            'icons': ['icon', 'icons'],
            'images': ['image', 'photo', 'picture'],
            'search': ['search', 'filter'],
            'dropdown': ['dropdown', 'select'],
            'tabs': ['tab', 'tabs'],
            'tooltip': ['tooltip', 'hover'],
            'loading': ['loading', 'spinner', 'loader']
        }
        
        for feature, keywords in feature_keywords.items():
            if any(keyword in prompt for keyword in keywords):
                features.append(feature)
        
        return features
    
    def _extract_styling_preferences(self, prompt: str) -> Dict:
        """Extract styling preferences"""
        styling = {
            'framework': 'tailwind',  # default
            'rounded': 'rounded' in prompt or 'curved' in prompt,
            'shadow': 'shadow' in prompt,
            'border': 'border' in prompt or 'outlined' in prompt,
            'gradient': 'gradient' in prompt
        }
        
        if 'bootstrap' in prompt:
            styling['framework'] = 'bootstrap'
        elif 'material' in prompt:
            styling['framework'] = 'material'
        
        return styling
    
    def _determine_framework(self, requirements: Dict) -> str:
        """Determine the best framework for the component"""
        # Default to React with TypeScript
        return 'react'
    
    def _determine_styling(self, requirements: Dict) -> str:
        """Determine styling approach"""
        return requirements['styling']['framework']
    
    def _generate_component(self, requirements: Dict, framework: str, styling: str) -> str:
        """Generate the actual component code"""
        
        if requirements['component_type'] == 'login_form':
            return self._generate_login_form(requirements, framework, styling)
        elif requirements['component_type'] == 'signup_form':
            return self._generate_signup_form(requirements, framework, styling)
        elif requirements['component_type'] == 'card':
            return self._generate_card(requirements, framework, styling)
        elif requirements['component_type'] == 'navigation':
            return self._generate_navigation(requirements, framework, styling)
        elif requirements['component_type'] == 'modal':
            return self._generate_modal(requirements, framework, styling)
        elif requirements['component_type'] == 'button':
            return self._generate_button(requirements, framework, styling)
        else:
            return self._generate_generic_component(requirements, framework, styling)
    
    def _generate_login_form(self, requirements: Dict, framework: str, styling: str) -> str:
        """Generate a login form component"""
        colors = requirements['colors']
        is_dark = colors['background'] == '#1f2937'
        has_glow = colors['accent'] == '#00ffff'
        
        glow_classes = 'shadow-lg shadow-cyan-500/50' if has_glow else ''
        bg_class = 'bg-gray-900' if is_dark else 'bg-white'
        text_class = 'text-white' if is_dark else 'text-gray-900'
        
        return f'''import React, {{ useState }} from 'react';

interface LoginFormProps {{
  onLogin?: (email: string, password: string) => void;
  className?: string;
}}

const {requirements['component_name']}: React.FC<LoginFormProps> = ({{ onLogin, className = '' }}) => {{
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState<{{email?: string, password?: string}}>({{}});

  const handleSubmit = async (e: React.FormEvent) => {{
    e.preventDefault();
    setIsLoading(true);
    setErrors({{}});
    
    // Basic validation
    const newErrors: {{email?: string, password?: string}} = {{}};
    if (!email) newErrors.email = 'Email is required';
    if (!password) newErrors.password = 'Password is required';
    
    if (Object.keys(newErrors).length > 0) {{
      setErrors(newErrors);
      setIsLoading(false);
      return;
    }}
    
    try {{
      if (onLogin) {{
        await onLogin(email, password);
      }}
    }} catch (error) {{
      console.error('Login failed:', error);
    }} finally {{
      setIsLoading(false);
    }}
  }};

  return (
    <div className={{`min-h-screen flex items-center justify-center {bg_class.replace('bg-', 'bg-gradient-to-br from-')} py-12 px-4 sm:px-6 lg:px-8 ${{className}}`}}>
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <h2 className="mt-6 text-3xl font-extrabold {text_class}">
            Sign in to your account
          </h2>
          <p className="mt-2 text-sm text-gray-600">
            Enter your credentials below
          </p>
        </div>
        <form className="mt-8 space-y-6" onSubmit={{handleSubmit}}>
          <div className="rounded-md shadow-sm -space-y-px">
            <div>
              <label htmlFor="email-address" className="sr-only">
                Email address
              </label>
              <input
                id="email-address"
                name="email"
                type="email"
                autoComplete="email"
                required
                className={{`appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 {text_class} rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm {bg_class} ${{has_glow ? 'focus:shadow-lg focus:shadow-cyan-500/50' : ''}}`}}
                placeholder="Email address"
                value={{email}}
                onChange={{(e) => setEmail(e.target.value)}}
              />
              {{errors.email && <p className="mt-1 text-sm text-red-600">{{errors.email}}</p>}}
            </div>
            <div>
              <label htmlFor="password" className="sr-only">
                Password
              </label>
              <input
                id="password"
                name="password"
                type="password"
                autoComplete="current-password"
                required
                className={{`appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 {text_class} rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm {bg_class} ${{has_glow ? 'focus:shadow-lg focus:shadow-cyan-500/50' : ''}}`}}
                placeholder="Password"
                value={{password}}
                onChange={{(e) => setPassword(e.target.value)}}
              />
              {{errors.password && <p className="mt-1 text-sm text-red-600">{{errors.password}}</p>}}
            </div>
          </div>

          <div className="flex items-center justify-between">
            <div className="flex items-center">
              <input
                id="remember-me"
                name="remember-me"
                type="checkbox"
                className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
              />
              <label htmlFor="remember-me" className="ml-2 block text-sm {text_class}">
                Remember me
              </label>
            </div>

            <div className="text-sm">
              <a href="#" className="font-medium text-indigo-600 hover:text-indigo-500">
                Forgot your password?
              </a>
            </div>
          </div>

          <div>
            <button
              type="submit"
              disabled={{isLoading}}
              className={{`group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed ${{has_glow ? 'shadow-lg shadow-indigo-500/50 hover:shadow-indigo-400/50' : ''}} ${{glow_classes}}`}}
            >
              {{isLoading ? (
                <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              ) : null}}
              {{isLoading ? 'Signing in...' : 'Sign in'}}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}};

export default {requirements['component_name']};
'''

    def _generate_card(self, requirements: Dict, framework: str, styling: str) -> str:
        """Generate a card component"""
        colors = requirements['colors']
        is_dark = colors['background'] == '#1f2937'
        
        return f'''import React from 'react';

interface CardProps {{
  title?: string;
  subtitle?: string;
  content?: React.ReactNode;
  image?: string;
  actions?: React.ReactNode;
  className?: string;
}}

const {requirements['component_name']}: React.FC<CardProps> = ({{
  title,
  subtitle,
  content,
  image,
  actions,
  className = ''
}}) => {{
  return (
    <div className={{`{'bg-gray-800' if is_dark else 'bg-white'} rounded-lg shadow-md overflow-hidden ${{className}}`}}>
      {{image && (
        <img 
          src={{image}} 
          alt={{title || 'Card image'}}
          className="w-full h-48 object-cover"
        />
      )}}
      <div className="p-6">
        {{title && (
          <h3 className="text-lg font-semibold {'text-white' if is_dark else 'text-gray-900'} mb-2">
            {{title}}
          </h3>
        )}}
        {{subtitle && (
          <p className="text-sm {'text-gray-400' if is_dark else 'text-gray-600'} mb-4">
            {{subtitle}}
          </p>
        )}}
        {{content && (
          <div className="{'text-gray-300' if is_dark else 'text-gray-700'} mb-4">
            {{content}}
          </div>
        )}}
        {{actions && (
          <div className="flex justify-end space-x-2">
            {{actions}}
          </div>
        )}}
      </div>
    </div>
  );
}};

export default {requirements['component_name']};
'''

    def _generate_generic_component(self, requirements: Dict, framework: str, styling: str) -> str:
        """Generate a generic component"""
        return f'''import React from 'react';

interface {requirements['component_name']}Props {{
  className?: string;
  children?: React.ReactNode;
}}

const {requirements['component_name']}: React.FC<{requirements['component_name']}Props> = ({{ className = '', children }}) => {{
  return (
    <div className={{`p-4 rounded-lg shadow-md bg-white ${{className}}`}}>
      <h2 className="text-xl font-bold mb-4">
        {requirements['component_name'].replace('Component', '')}
      </h2>
      <div>
        {{children || <p>Custom component content goes here.</p>}}
      </div>
    </div>
  );
}};

export default {requirements['component_name']};
'''

    def _save_code(self, code: str, component_name: str, framework: str) -> str:
        """Save generated code to file"""
        # Ensure output directory exists
        self.output_dir.mkdir(parents=True, exist_ok=True)
        
        # Determine file extension
        ext = '.tsx' if framework == 'react' else '.vue' if framework == 'vue' else '.ts'
        
        # Create filename
        filename = f"{component_name}{ext}"
        file_path = self.output_dir / filename
        
        # Write code to file
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(code)
        
        return str(file_path)
    
    def _get_react_templates(self) -> Dict:
        """Get React component templates"""
        return {}  # Implemented above in generate methods
    
    def _get_vue_templates(self) -> Dict:
        """Get Vue component templates"""
        return {}  # Could be implemented for Vue support
    
    def _get_angular_templates(self) -> Dict:
        """Get Angular component templates"""
        return {}  # Could be implemented for Angular support
    
    def _get_tailwind_classes(self) -> Dict:
        """Get Tailwind CSS class mappings"""
        return {
            'colors': {
                'primary': 'bg-blue-500',
                'secondary': 'bg-gray-500',
                'success': 'bg-green-500',
                'danger': 'bg-red-500',
                'warning': 'bg-yellow-500',
                'info': 'bg-blue-400'
            },
            'spacing': {
                'xs': 'p-2',
                'sm': 'p-4',
                'md': 'p-6',
                'lg': 'p-8',
                'xl': 'p-10'
            }
        }
    
    def _get_bootstrap_classes(self) -> Dict:
        """Get Bootstrap class mappings"""
        return {}  # Could be implemented
    
    def _get_material_classes(self) -> Dict:
        """Get Material Design class mappings"""
        return {}  # Could be implemented