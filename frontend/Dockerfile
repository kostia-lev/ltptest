# Use the official Node.js image for building
FROM node:18 as builder

WORKDIR /app

# Copy package.json and package-lock.json files
COPY package.json package-lock.json ./

# Install dependencies using npm
RUN npm install

# Copy source code
COPY . .

# Build the React app
RUN npm run build

# Use Nginx for serving the built app
FROM nginx:alpine

# Copy the built React app to the Nginx HTML directory
COPY --from=builder /app/dist /usr/share/nginx/html

# Expose the port used by Nginx
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
