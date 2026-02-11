FROM node:16-bullseye-slim
MAINTAINER OpsTree Solutions

# Set working directory
WORKDIR /app

# Install git (using a modern base image where repositories work)
RUN apt-get update -y && \
    apt-get install git -y

# Copy package files and install dependencies
COPY package*.json ./
RUN npm install

RUN npm run build

RUN yarn global add serve

# Copy the rest of the code
COPY . .

EXPOSE 80
ENTRYPOINT ["serve", "-s", "build"]
