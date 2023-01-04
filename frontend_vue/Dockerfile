FROM node:alpine as build-stage

ARG FRONTEND_LOCATION
ARG BACKEND_URL

WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN echo VUE_APP_BACKEND_URL="${BACKEND_URL}">.env
RUN if [ "${FRONTEND_LOCATION}" != "" ]; then \
 export PUBLIC_PATH="${FRONTEND_LOCATION}"; \
 export BASE_URL="${FRONTEND_LOCATION}"; \
 export VUE_APP_BACKEND_URL="${BACKEND_URL}"; \
 npm run build -- "${FRONTEND_LOCATION}"; \
else \
 npm run build; \
fi

FROM nginx:alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
