import React from 'react';
import LearnWord from './pages/LearnWord';
import Schedule from './pages/Schedule';
import ThemeProvider from './components/ThemeProvider';
import { WordProvider } from './context/WordContext';
import './styles/global.css';

const App = () => {
  return (
    <ThemeProvider>
      <WordProvider>
        <div className="grid gap-8">
          <LearnWord />
          <Schedule />
        </div>
      </WordProvider>
    </ThemeProvider>
  );
};

export default App;
