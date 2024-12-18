FROM node:latest

ENV PORT=5501

WORKDIR /Documents/learning-english/analysis

COPY ./analysis/frontend/package*.json ./
RUN npm install

COPY ./analysis/frontend/ ./

EXPOSE 5501

CMD ["npm", "start"]
