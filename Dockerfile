FROM node

WORKDIR /app
COPY . .

RUN npm i

ENTRYPOINT ["node", "index.js"]
