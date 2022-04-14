FROM node

WORKDIR /app
COPY ./src .

RUN npm i

ENTRYPOINT ["node", "index.js"]
