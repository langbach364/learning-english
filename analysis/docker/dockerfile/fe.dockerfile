FROM node:latest

WORKDIR /Documents/learning-english/analysis
COPY ./frontend/package*.json ./
RUN npm install

COPY ./frontend .
EXPOSE 5501
CMD ["npm", "start"]
