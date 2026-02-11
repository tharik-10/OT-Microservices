FROM node:16-bullseye-slim
MAINTAINER OpsTree Solutions

WORKDIR /app

RUN apt-get update -y && \
    apt-get install -y git

# Copy package files first (better layer caching)
COPY package*.json ./
RUN npm install

# Now copy full source code
COPY . .

# Now build (after source exists)
RUN npm run build

RUN yarn global add serve

EXPOSE 3000

ENTRYPOINT ["serve", "-s", "build", "-l", "3000"]
