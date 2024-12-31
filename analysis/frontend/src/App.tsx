import React, { useEffect, useState, useCallback } from "react";
import { APIService } from "./services/api";
import { AnswerData } from "./types/dictionary";
import ConnectionStatus from "./components/ConnectionStatus";
import WordInput from "./components/WordInput";
import ChatDisplay from "./components/ChatDisplay";
import Loading from "./components/Loading";
import { Login } from "./components/Login";
import styled from "@emotion/styled";
import { fadeIn } from "./styles/animation";

const AppContainer = styled.div`
  animation: ${fadeIn} 0.5s ease-out;
  min-height: 100vh;
  background-color: rgb(249 250 251);
  padding: 2rem 0;
`;

const ContentWrapper = styled.div`
  max-width: 64rem;
  margin: 0 auto;
  padding: 0 1rem;
`;

const Title = styled.h1`
  color: #2563eb;
  font-size: 2.25rem;
  font-weight: 700;
  text-align: center;
  margin-bottom: 2rem;
`;

function App() {
  const [word, setWord] = useState("");
  const [definitions, setDefinitions] = useState<AnswerData>({
    detail: {},
    structure: "Sentence",
  });
  const [messageType, setMessageType] = useState<string>("dictionary");
  const [isLoading, setIsLoading] = useState(false);
  const [isConnected, setIsConnected] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const handleWebSocketMessage = useCallback(
    (data: AnswerData, type: string) => {
      if (data) {
        setDefinitions({
          detail: data.detail || {},
          structure: data.structure || "Sentence",
        });
        setMessageType(type || "dictionary");
        setIsLoading(false);
      }
    },
    []
  );

  const verifyAuthentication = async () => {
    try {
      const response = await APIService.verifyToken();
      if (!response.success) {
        localStorage.removeItem('token');
        setIsAuthenticated(false);
      }
    } catch (error) {
      console.error('Verification failed:', error);
      localStorage.removeItem('token');
      setIsAuthenticated(false);
    }
  };

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      setIsAuthenticated(true);
      verifyAuthentication();
      const cleanup = APIService.connectWebSocket(
        handleWebSocketMessage,
        setIsConnected
      );
      return cleanup;
    }
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!word.trim() || !isAuthenticated) return;

    setIsLoading(true);
    try {
      await APIService.searchWord(word);
      setWord("");
    } catch (error) {
      console.error("Request error:", error);
      setIsLoading(false);
    }
  };

  const handleLoginSuccess = useCallback(() => {
    setIsAuthenticated(true);
    verifyAuthentication(); // Verify sau khi đăng nhập thành công
  }, []);

  return (
    <AppContainer>
      {!isAuthenticated ? (
        <Login onLoginSuccess={handleLoginSuccess} />
      ) : (
        <>
          {isLoading && <Loading />}
          <ConnectionStatus isConnected={isConnected} />

          <ContentWrapper>
            <Title>Phân tích từ hoặc câu tiếng Anh</Title>

            <WordInput
              word={word}
              setWord={setWord}
              onSubmit={handleSubmit}
              isLoading={isLoading}
              data={definitions.detail}
            />

            <ChatDisplay data={definitions} type={messageType} />
          </ContentWrapper>
        </>
      )}
    </AppContainer>
  );
}

export default App;
