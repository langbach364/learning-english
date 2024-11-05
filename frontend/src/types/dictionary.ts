export interface GrammarAnalysis {
  detail: {
    [key: string]: string[];
  };
  structure: "WordClass" | "Sentence";
}

export interface DictionaryExample {
  EN: string[];
  VI: string[];
}

export interface DictionaryDefinition {
  [definition: string]: DictionaryExample;
}

export interface DictionaryWordType {
  [wordType: string]: DictionaryDefinition;
}

export interface GrammarDetail {
  [key: string]: string[];
}

export interface AnswerData {
  detail: DictionaryWordType | GrammarDetail;
  structure: "WordClass" | "Sentence";
}

export interface WebSocketMessage {
  type: string;
  data: AnswerData;
}
