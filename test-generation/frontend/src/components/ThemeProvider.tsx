import React from 'react';
import { ThemeProvider as StyledThemeProvider } from 'styled-components';

const theme = {
  colors: {
    primary: '#2563eb',
    secondary: '#4f46e5',
    background: '#f3f4f6',
    text: '#1f2937'
  },
  fonts: {
    body: 'Inter, system-ui, sans-serif',
    heading: 'Poppins, sans-serif'
  }
};

const ThemeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <StyledThemeProvider theme={theme}>
      {children}
    </StyledThemeProvider>
  );
};
export default ThemeProvider;
