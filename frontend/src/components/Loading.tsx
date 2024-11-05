import React from 'react';
import styled from '@emotion/styled';
import { keyframes } from '@emotion/react';

const spin = keyframes`
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
`;

const pulse = keyframes`
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
`;

const LoadingContainer = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  z-index: 1000;
`;

const Spinner = styled.div`
  width: 50px;
  height: 50px;
  border: 4px solid #EFF6FF;
  border-top: 4px solid #3B82F6;
  border-radius: 50%;
  animation: ${spin} 1s linear infinite;
`;

const LoadingText = styled.div`
  margin-top: 1rem;
  color: #3B82F6;
  font-weight: 500;
  animation: ${pulse} 1.5s ease-in-out infinite;
`;

const Loading: React.FC = () => {
  return (
    <LoadingContainer>
      <Spinner />
      <LoadingText>AI đang trong quá trình xử lý...</LoadingText>
    </LoadingContainer>
  );
};

export default Loading;
