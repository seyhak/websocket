FROM node:22-alpine

WORKDIR /usr/src/app

COPY ./UI .

RUN ls -a
RUN npm i

CMD ["npm run", "dev"]
