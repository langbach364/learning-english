import React from 'react';
import LearnWord from './pages/LearnWord';
import Schedule from './pages/Schedule';
import Statistics from './pages/Statistics';
import ThemeProvider from './components/ThemeProvider';
import { WordProvider } from './context/WordContext';
import './styles/global.css';

const App: React.FC = () => {
  return (
    <ThemeProvider>
      <WordProvider>
        <div className="min-h-screen bg-gray-50 p-4">
          <div className="max-w-7xl mx-auto">
            <div className="grid grid-cols-1 gap-8">
              <LearnWord />
              <Schedule />
              <Statistics />
            </div>
          </div>
        </div>
      </WordProvider>
    </ThemeProvider>
  );
};

export default App;
