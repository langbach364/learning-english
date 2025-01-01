import React, { useState, useEffect } from 'react';
import { APIService } from '../services/api';
import { API_CONFIG } from '../constants/config';
import styled from '@emotion/styled';
import { fadeIn, slideDown } from '../styles/animation';

const LoginContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background-color: rgb(249 250 251);
  padding: 1rem;
  animation: ${fadeIn} 0.5s ease-out;

  @media (min-width: 640px) {
    padding: 2rem;
  }
`;

const LoginBox = styled.div`
  background: white;
  padding: 1.5rem;
  border-radius: 0.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 320px;
  animation: ${slideDown} 0.5s ease-out;

  @media (min-width: 640px) {
    padding: 2rem;
    max-width: 400px;
    border-radius: 0.75rem;
  }

  @media (min-width: 768px) {
    padding: 2.5rem;
    max-width: 450px;
  }
`;

const Title = styled.h1`
  color: #2563eb;
  font-size: var(--font-size-xl);
  text-align: center;
  margin-bottom: 1.5rem;
  font-weight: 600;

  @media (min-width: 640px) {
    font-size: var(--font-size-2xl);
    margin-bottom: 2rem;
  }

  @media (min-width: 768px) {
    font-size: var(--font-size-3xl);
  }
`;

const Button = styled.button<{ isLoading?: boolean; isLogout?: boolean }>`
  width: 100%;
  padding: 0.75rem;
  background-color: ${props => {
    if (props.isLogout) return '#dc2626';
    return props.isLoading ? '#93c5fd' : '#2563eb';
  }};
  color: white;
  border: none;
  border-radius: 0.375rem;
  font-size: var(--font-size-base);
  cursor: ${props => props.isLoading ? 'not-allowed' : 'pointer'};
  transition: all 0.2s ease;
  margin-top: ${props => props.isLogout ? '1rem' : '0'};

  &:hover {
    background-color: ${props => {
      if (props.isLogout) return '#b91c1c';
      return props.isLoading ? '#93c5fd' : '#1d4ed8';
    }};
    transform: translateY(-2px);
  }

  &:active {
    transform: translateY(1px);
  }

  @media (min-width: 640px) {
    padding: 1rem;
    font-size: var(--font-size-lg);
  }
`;

const Message = styled.div<{ type: 'error' | 'status' }>`
  color: ${props => props.type === 'error' ? '#dc2626' : '#059669'};
  background-color: ${props => props.type === 'error' ? '#fee2e2' : '#d1fae5'};
  padding: 0.75rem;
  border-radius: 0.375rem;
  margin-bottom: 1rem;
  text-align: center;
  font-size: var(--font-size-sm);

  @media (min-width: 640px) {
    padding: 1rem;
    font-size: var(--font-size-base);
  }
`;

interface LoginProps {
  onLoginSuccess: () => void;
}

export const Login: React.FC<LoginProps> = ({ onLoginSuccess }) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [status, setStatus] = useState<string | null>(null);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      handleAutoLogin();
    }
  }, []);

  const handleAutoLogin = async () => {
    try {
      setStatus('Đang kiểm tra phiên đăng nhập...');
      onLoginSuccess();
    } catch (err) {
      localStorage.removeItem('token');
      setError('Phiên đăng nhập đã hết hạn');
    } finally {
      setStatus(null);
    }
  };

  const handleLogin = async () => {
    try {
      setLoading(true);
      setError(null);
      setStatus('Đang đăng nhập...');
  
      const username = API_CONFIG.AUTH.CREDENTIALS.USERNAME;
      const password = API_CONFIG.AUTH.CREDENTIALS.PASSWORD;

      if (!username || !password) {
        throw new Error('Tên đăng nhập hoặc mật khẩu không được để trống');
      }

      await APIService.login({
        username,
        password
      });
  
      setStatus('Đăng nhập thành công!');
      
      setTimeout(() => {
        onLoginSuccess();
      }, 1000);
  
    } catch (err) {
      console.error('Lỗi đăng nhập:', err);
      setError(err instanceof Error ? err.message : 'Đăng nhập thất bại');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    APIService.logout();
    setStatus(null);
    setError(null);
    window.location.reload();
  };

  return (
    <LoginContainer>
      <LoginBox>
        <Title>Đăng nhập hệ thống</Title>

        {error && <Message type="error">{error}</Message>}
        {status && <Message type="status">{status}</Message>}

        <Button 
          onClick={handleLogin} 
          disabled={loading}
          isLoading={loading}
        >
          {loading ? 'Đang xử lý...' : 'Đăng nhập'}
        </Button>

        {localStorage.getItem('token') && (
          <Button 
            onClick={handleLogout}
            isLogout
          >
            Đăng xuất
          </Button>
        )}
      </LoginBox>
    </LoginContainer>
  );
}
