import 'vue-universal-modal/dist/index.css'
import VueUniversalModal from 'vue-universal-modal'
import { createApp } from 'vue'
import App from './App.vue'
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { library } from "@fortawesome/fontawesome-svg-core";
import { faTrashCan } from "@fortawesome/free-solid-svg-icons";

library.add(faTrashCan);

const app = createApp(App);
app.use(VueUniversalModal, {
  teleportTarget: '#modals'
})
app.component("font-awesome-icon", FontAwesomeIcon);
app.mount('#app');
