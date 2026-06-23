FROM node:18-alpine

WORKDIR /app

# Copiar configuración de dependencias y realizar la instalación
COPY package.json ./
RUN npm install

# Copiar todo el código del frontend
COPY . .

EXPOSE 3000

# Levantar el servidor de desarrollo de Vite exponiendo el host
CMD ["npm", "run", "dev", "--", "--host"]
