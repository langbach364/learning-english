FROM node:latest

ENV PORT=5501
ENV NODE_ENV=production
ENV NODE_OPTIONS=--openssl-legacy-provider

WORKDIR /Documents/learning-english/analysis

COPY ./analysis/frontend/package*.json ./
RUN npm install --production

COPY ./analysis/frontend/ ./
RUN npm run build
RUN npm install -g serve

EXPOSE 5501

CMD ["serve", "-s", "build", "-l", "5501"]