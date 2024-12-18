FROM node:latest

WORKDIR /Documents/learning-english/analysis
COPY ./analysis/frontend/ ./

ENV PORT=5501
EXPOSE 5501
CMD ["npm", "start"]
