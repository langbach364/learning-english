export interface Word {
  id: number;
  word: string;
  wrongCount: number;
}

export interface StatisticsPayload {
  date: string;
  range: "DAY" | "WEEK" | "MONTH";
}
