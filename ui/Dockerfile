# ui/Dockerfile
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package.json and install dependencies
#COPY package.json package-lock.json ./
COPY package*.json ./

RUN npm install

# Copy the rest of the application code
COPY . .

# Build the application
RUN npm run build

# Install serve to serve the static files
RUN npm install -g serve

# Expose the port the app runs on
EXPOSE 3000

# Command to serve the application
CMD ["serve", "-s", "dist", "-l", "3000"]