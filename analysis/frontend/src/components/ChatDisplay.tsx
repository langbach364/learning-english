import React from "react";
import { AnswerData, DictionaryWordType, GrammarDetail } from "../types/dictionary";
import DictionaryDisplay from "./DictionaryDisplay";
import GrammarDisplay from "./GrammarDisplay";

interface ChatDisplayProps {
  data: AnswerData;
  type: string;
}

const ChatDisplay: React.FC<ChatDisplayProps> = ({ data, type }) => {
  if (!data?.detail) return null;

  const displayData = {
    detail: data.detail,
    structure: data.structure
  };

  const renderContent = () => {
    if (displayData.structure === "Sentence") {
      return <GrammarDisplay 
        content={{
          detail: displayData.detail as GrammarDetail,
          structure: displayData.structure
        }} 
      />;
    } else if (displayData.structure === "WordClass") {
      return <DictionaryDisplay 
        content={displayData.detail as DictionaryWordType}
        structure={displayData.structure}
      />;
    }
    return null;
  };

  return <div className="mt-8">{renderContent()}</div>;
};

export default ChatDisplay;
