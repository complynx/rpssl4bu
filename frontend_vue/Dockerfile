FROM node:alpine as build-stage

ARG FRONTEND_LOCATION

WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN if [ "${FRONTEND_LOCATION}" != "" ]; then \
 export PUBLIC_PATH="${FRONTEND_LOCATION}"; \
 export BASE_URL="${FRONTEND_LOCATION}"; \
 npm run build -- "${FRONTEND_LOCATION}"; \
else \
 npm run build; \
fi

FROM nginx:alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
