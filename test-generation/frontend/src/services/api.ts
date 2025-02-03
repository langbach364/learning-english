import axios from 'axios';
import { Word } from '../types';

const REST_API = 'http://localhost:8081';

export const api = axios.create({
  baseURL: REST_API,
  headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
  }
});

export const learnWord = async (count: number) => {
  try {
      const response = await api.post('/learn_word', { count });
      return response.data;
  } catch (error) {
      console.error('Error fetching words:', error);
      throw error;
  }
};

export const createSchedule = async (words: Word[]) => {
  try {
    const response = await api.post('/create_schedule', words);
    return response.data;
  } catch (error) {
    console.error('Error creating schedule:', error);
    throw error;
  }
};

export const reviseWord = async () => {
  return await api.post('/revise_word');
};

export const getStatistics = async () => {
  return await api.post('/get_statistics');
};