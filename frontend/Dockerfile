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

# Copy the rest of the code
COPY . .

EXPOSE 3000
CMD ["npm", "start"]
