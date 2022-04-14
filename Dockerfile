FROM node

WORKDIR /app
COPY . .

RUN npm i

ENTRYPOINT ["node", "src/index.js"]
