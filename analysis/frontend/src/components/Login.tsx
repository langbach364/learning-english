import React, { useState, useEffect } from 'react';
import { APIService } from '../services/api';
import { API_CONFIG } from '../constants/config';
import styled from '@emotion/styled';

const LoginContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background-color: rgb(249 250 251);
`;

const LoginBox = styled.div`
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
`;

const Title = styled.h1`
  color: #2563eb;
  font-size: 1.5rem;
  text-align: center;
  margin-bottom: 2rem;
`;

const Button = styled.button<{ isLoading?: boolean }>`
  width: 100%;
  padding: 0.75rem;
  background-color: ${props => props.isLoading ? '#93c5fd' : '#2563eb'};
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: ${props => props.isLoading ? 'not-allowed' : 'pointer'};
  transition: background-color 0.2s;

  &:hover {
    background-color: ${props => props.isLoading ? '#93c5fd' : '#1d4ed8'};
  }
`;

const ErrorMessage = styled.div`
  color: #dc2626;
  background-color: #fee2e2;
  padding: 0.75rem;
  border-radius: 4px;
  margin-bottom: 1rem;
  text-align: center;
`;

const StatusMessage = styled.div`
  color: #059669;
  background-color: #d1fae5;
  padding: 0.75rem;
  border-radius: 4px;
  margin-bottom: 1rem;
  text-align: center;
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

        {error && <ErrorMessage>{error}</ErrorMessage>}
        {status && <StatusMessage>{status}</StatusMessage>}

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
            style={{ marginTop: '1rem', backgroundColor: '#dc2626' }}
          >
            Đăng xuất
          </Button>
        )}
      </LoginBox>
    </LoginContainer>
  );
}