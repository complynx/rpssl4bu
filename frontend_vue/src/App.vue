<template>
  <nav>
    <ul>
      <li v-if="!p2pMode">
        <a @click="isShowStartP2P=true">Play P2P</a>
      </li>
      <li v-if="p2pMode">
        <a @click="copyInviteLinkToClipboard">Copy invitation link</a>
      </li>
      <li v-if="p2pMode">
        <a @click="cancelP2P">Exit P2P mode</a>
      </li>
      <li v-if="!p2pMode">
        <a v-if="globalResults" @click="toggleResults">Only local results</a>
        <a v-if="!globalResults" @click="toggleResults">Global results</a>
      </li>
      <li>
        <a @click="openChangeBackendModal">Change backend server</a>
      </li>
    </ul>
  </nav>
  <main>
    <div class="game">
      <div class="player player-left">
        <h1 v-if="!p2pMode">You</h1>
        <h1 v-if="p2pMode">
          {{ leftPlayerName }}
          <span v-if="leftPlayerName == ''" class="waiting-player">waiting...</span>
          <span v-if="!yourSideIsRight" class="player-you-mark"> (you)</span>
        </h1>
        <template v-if="p2pMode">
          <div v-if="leftPlayerChoiceID != 0" class="weapon" :class="'weapon-'+leftPlayerChoiceID">
            {{ leftPlayerChoice }}
          </div>
        </template>
        <template v-if="!p2pMode">
          <div v-if="yourChoiceId != 0" class="weapon" :class="'weapon-'+yourChoiceId">
            {{ yourChoice }}
          </div>
        </template>
      </div>
      <div class="weapons">
        <h1>Choose your weapon</h1>
        <div class="weapons-choices">
          <button class="weapon" v-for="weapon in weapons" :key="weapon.id" :class="'weapon-'+weapon.id" @click="makeChoice(weapon.id)">
            {{ weapon.name }}
          </button>
        </div>
      </div>
      <div class="player player-right">
        <h1 v-if="!p2pMode">Computer</h1>
        <h1 v-if="p2pMode">
          {{ rightPlayerName }}
          <span v-if="rightPlayerName == ''" class="waiting-player">waiting...</span>
          <span v-if="yourSideIsRight" class="player-you-mark"> (you)</span>
        </h1>
        <template v-if="p2pMode">
          <div v-if="rightPlayerChoiceID != 0" class="weapon" :class="'weapon-'+rightPlayerChoiceID">
            {{ rightPlayerChoice }}
          </div>
        </template>
        <template v-if="!p2pMode">
          <div v-if="computerChoiceId != 0" class="weapon" :class="'weapon-'+computerChoiceId">
            {{ computerChoice }}
          </div>
        </template>
      </div>
    </div>
    <div class="board">
      <h1 v-if="!globalResults || p2pMode">Your last scores:</h1>
      <h1 v-if="globalResults && !p2pMode">Last global scores:</h1>
      <div class="scores">
        <template v-if="globalResults && !p2pMode">
          <div class="score" v-for="(score, index) in scores" :key="index">
            {{ score }}
          </div>
        </template>
        <template v-if="!globalResults && !p2pMode">
          <div class="score" v-for="(score, index) in localScores" :key="index">
            {{ score }}
          </div>
        </template>
        <template v-if="p2pMode">
          <div class="score" v-for="(score, index) in p2pScores" :key="index">
            {{ score }}
          </div>
        </template>
      </div>
      <button @click="clearScores()">
        <font-awesome-icon icon="fa-solid fa-trash-can"/>
      </button>
      
    </div>
  </main>
  <Modal
    v-model="isShowResult"
    :close="closeResultModal"
  >
    <div class="modal result">
      <h1 v-if="result === 'win'">You win!</h1>
      <h1 v-if="result === 'lose'">You've lost</h1>
      <h1 v-if="result === 'tie'">It's a tie</h1>
      <div class="content">
        <template v-if="resultRepresentation === null && !p2pMode">
          <div>Computer chose {{computerChoice}}</div>
          <div>You chose {{yourChoice}}</div>
        </template>
        <template v-if="resultRepresentation === null && p2pMode && yourSideIsRight">
          <div>{{leftPlayerName}} chose {{leftPlayerChoice}}</div>
          <div>You chose {{rightPlayerChoice}}</div>
        </template>
        <template v-if="resultRepresentation === null && p2pMode && !yourSideIsRight">
          <div>You chose {{leftPlayerChoice}}</div>
          <div>{{rightPlayerName}} chose {{rightPlayerChoice}}</div>
        </template>
        <template v-if="resultRepresentation !== null">
          <div>{{resultRepresentation.text}}</div>
          <img :src="'./images/'+resultRepresentation.img" />
        </template>
      </div>
    </div>
  </Modal>
  <Modal
    v-model="isShowChangeBackendServer"
    :close="closeChangeBackendServerModal"
  >
    <div class="modal backend-server">
      <h1>Change backend server</h1>
      <div class="content">
        Enter the new backend server address:
        <input v-model="backendServer" type="text">
        Note, that some servers may lack P2P or score board functionality
      </div>
      <div class="footer">
        <button @click="changeBackendServer">Save</button>
        <button @click="closeChangeBackendServerModal">
          Cancel
        </button>
      </div>
    </div>
  </Modal>
  <Modal
    v-model="isShowStartP2P"
    :close="cancelP2P"
  >
    <div class="modal start-p2p">
      <h1>Start P2P game</h1>
      <div class="content">
        <button @click="createP2P">Create P2P game</button>
        <div>Or post the invitation link:</div>
        <input @input="inviteChanged" type="text">
        <div class="input-check-error" v-if="showStartError==true">Game not found or full.</div>
      </div>
    </div>
  </Modal>
  <Modal
    v-model="isShowJoinP2P"
    :close="cancelP2P"
  >
    <div class="modal join-p2p">
      <h1>Join P2P game</h1>
      <div class="content">
        <div>Copy and send this to your friend:</div>
        <input :value="p2pInviteLink" readonly type="text" @click="copyInviteLinkToClipboard" />
        <div>Enter your name:</div>
        <input v-model="yourName" ref="joinP2PNameInput" type="text" pattern="^[a-zA-Z ]{0,20}$">
      </div>
      <div class="footer">
        <button @click="joinP2P">Join</button>
      </div>
    </div>
  </Modal>
</template>

<script>
import axios from 'axios';

function hexToB64(str) {
  return btoa(String.fromCharCode.apply(null, str.match(/\w{2}/g).map(function(a) {
      return parseInt(a, 16);
  }))).replace('/','-').replace('+','_');
}
function b64ToHex(str) {
  return atob(str.replace('-','/').replace('_','+')).split('').map(function(c) {
      return ('00' + c.charCodeAt(0).toString(16)).slice(-2);
  }).join('');
}

export default {
  name: 'App',
  components: {
  },
  data() {
    return {
      backendServer: '/backend/', // prod
      // backendServer: 'http://localhost:80/backend/', // for easy testing
      // backendServer: 'https://codechallenge.boohma.com/',
      isShowChangeBackendServer: false,
      isShowResult: false,
      isShowStartP2P: false,
      showStartError: false,
      isShowJoinP2P: false,
      p2pMode: false,
      p2pID: "",
      p2pInviteLink: "",
      yourName: "",
      computerChoice: '',
      yourChoice: '',
      computerChoiceId: 0,
      globalResults: true,
      yourChoiceId: 0,
      result: '',
      resultRepresentation: null,
      leftPlayerName: "",
      rightPlayerName: "",
      leftPlayerChoice: "",
      rightPlayerChoice: "",
      leftPlayerChoiceID: 0,
      rightPlayerChoiceID: 0,
      yourSideIsRight: false,
      scores: [],
      localScores: [],
      p2pScores: [],
      weapons: [
        { id: 1, name: 'Rock' },
        { id: 2, name: 'Paper' },
        { id: 3, name: 'Scissors' },
        { id: 4, name: 'Lizard' },
        { id: 5, name: 'Spock' }
      ],
      weaponDict: {}
    }
  },
  created() {
    this.populateWeapons();
    this.fetchWeapons();
    this.fetchScores();

    window.addEventListener("hashchange", this.hashProcess.bind(this));
    this.hashProcess();
    
    this.yourName = localStorage.getItem("p2p-player-name") || "";

    window.addEventListener("resize", function(){
      window?.resize_weapons?.apply?.(this, arguments);
    }, {
      'passive': true
    });
    window.resize_weapons = function(){
      let els = document.getElementsByClassName('weapons-choices');
      if(els.length < 1) return;
      let weapons = els[0];

      let h = weapons.clientHeight, w = weapons.clientWidth;
      let min = h < w ? h : w;
      let optimal_weapon_size = Math.floor(.3 * min);
      weapons.style.setProperty('--weapon-size', optimal_weapon_size + 'px');
    };
    // let it all render first
    requestAnimationFrame(window.resize_weapons);
  },
  methods: {
    getResultRepresentation(p1,p2) {
      if(p1>p2) [p2,p1]=[p1,p2];
      let key = (p1 + '_' + p2).toLowerCase();

      switch(key){
        case 'rock_rock':
          return {
            img: "rock_rock.png",
            text: "Now you have a pile of rocks in a tie!"
          };
        case 'paper_rock':
          return {
            img: "rock_paper.png",
            text: "A gift for naughty children?"
          };
        case 'rock_scissors':
          return {
            img: "rock_scissors.png",
            text: "Now you have some scrap metal. What do you want to do?"
          };
        case 'lizard_rock':
          return {
            img: "lizard_rock.png",
            text: "A rock landed on a planet somewhere near Mexico, archaeologists still find evidence of giant lizards from that era."
          };
        case 'scissors_scissors':
          return {
            img: "hair_salon.png",
            text: "Are you opening a hair salon?"
          };
        case 'lizard_lizard':
          return {
            img: "fancy_lizards.png",
            text: "You see two lizards in fancy ties. They might be tax agents!"
          };
        case 'spock_spock':
          return {
            img: "2spocks.png",
            text: "Clones war? Or is it the dark truth of how teleporters really work?"
          };
        case 'paper_paper':
          return {
            img: "thesis.png",
            text: "You both want to become PhD. Now it is time to defend your thesis."
          };
        case 'paper_scissors':
          return {
            img: "kirigami.png",
            text: "Achievement unlocked: Master of art of kirigami."
          };
        case 'paper_spock':
          return {
            img: "spock_taxes.png",
            text: "Even for Spock filling out taxes is too much."
          };
        case 'lizard_paper':
          return {
            img: "lizard_paper.png",
            text: "A lizard ate your homework. Why teachers don't want to believe you?"
          };
        case 'lizard_spock':
          return {
            img: "spock_lizard.png",
            text: "He didn't have time to turn around and see what was coming."
          };
        case 'scissors_spock':
          return {
            img: "laser_scissors.png",
            text: "Molten metal is very hot and bright."
          };
        case 'rock_spock':
          return {
            img: "spock_rock.png",
            text: "All the rocks are belong to us."
          };
        case 'lizard_scissors':
          return {
            img: "sl.png",
            text: "Never cut the wrong side of the lizard."
          };
      }
      return null;
    },
    async copyInviteLinkToClipboard() {
      try {
        await navigator.clipboard.writeText(window.location.href.split("#")[0] + "#" + this.p2pInviteLink);
      } catch (error) {
        console.error(error);
      }
    },
    joinP2P() {
      if(!this.$refs.joinP2PNameInput.validity.valid) return;
      this.isShowJoinP2P = false;

      let srv = new String(this.backendServer);
      if (srv.startsWith('https://')) {
        srv = srv.replace(/https:\/\//, 'wss://');
      } else {
        srv = srv.replace(/http:\/\//, 'ws://');
      }
      
      this.p2pSocket = new WebSocket(srv +
        'connect_p2p?g=' +
        encodeURIComponent(this.p2pID) +
        "&name=" +
        encodeURIComponent(this.yourName), "p2p");
      console.log(this.p2pSocket);

      this.bindSocket();
    },
    async checkP2PGame() {
      let id = this.p2pID;
      try {
        const response = await axios.get(this.backendServer + 'find_p2p?g=' + encodeURIComponent(id));
        if(this.p2pID != id) return;
        if(!response.data.is_full){
          this.showStartError = false;
          this.isShowJoinP2P = true;
          return;
        }
      }catch (err) {
        if(this.p2pID != id) return;
      }
      if(!this.isShowStartP2P) this.cancelP2P();
      else this.showStartError = true;
    },
    bindSocket() {
      if(!this.p2pSocket) return;
      this.p2pSocket.onopen = this.onSocketOpen.bind(this);
      this.p2pSocket.onerror = this.onSocketError.bind(this);
      this.p2pSocket.onclose = this.onSocketClose.bind(this);
      this.p2pSocket.onmessage = this.onSocketMessage.bind(this);
    },
    onSocketOpen(ev) {
      console.log(ev);
      this.p2pMode = true;
      this.p2pScores = [];
    },
    onSocketError(ev) {
      console.log(ev);
    },
    onSocketClose() {
      this.cancelP2P();
    },
    onSocketMessage(ev) {
      let data = JSON.parse(ev.data);
      console.log(data);
      this.leftPlayerName = data.state.left_player_name;
      this.rightPlayerName = data.state.right_player_name;
      this.yourSideIsRight = data.side != "left";
      this.leftPlayerChoice = data.state.left_player_choice.name;
      this.leftPlayerChoiceID = data.state.left_player_choice.id;
      this.rightPlayerChoice = data.state.right_player_choice.name;
      this.rightPlayerChoiceID = data.state.right_player_choice.id;
      if(data.state.result != "unknown") {
        this.result = data.state.result;
        this.resultRepresentation = this.getResultRepresentation(this.leftPlayerChoice, this.rightPlayerChoice);
        this.isShowResult = true;
        this.p2pScores.unshift(this.result);
        if(this.p2pScores.length>10) this.p2pScores = this.p2pScores.slice(0, 10);
      }
    },
    unbindSocket() {
      if(!this.p2pSocket) return;
      this.p2pSocket.onopen = null;
      this.p2pSocket.onerror = null;
      this.p2pSocket.onclose = null;
      this.p2pSocket.onmessage = null;
      this.p2pSocket.close();
      this.p2pSocket = null;
    },
    cancelP2P() {
      this.p2pMode = false;
      this.isShowStartP2P = false;
      this.isShowJoinP2P = false;
      this.showStartError = false;
      this.p2pID = "";
      this.p2pInviteLink = "";
      location.hash = "";
      this.unbindSocket();
    },
    async createP2P() {
      const response = await axios.post(this.backendServer + 'create_p2p');
      this.isShowStartP2P = false;
      this.p2pID = response.data;
      this.isShowJoinP2P = true;
      this.p2pInviteLink = hexToB64(this.p2pID).substring(0,11);
      location.hash = "#" + this.p2pInviteLink;
    },
    inviteChanged(ev){
      let val = ev.target.value;
      if(val.length == 11) return this.inviteProcess(val);
      if(val.indexOf("#")>=0) {
        return this.inviteProcess(val.substring(val.indexOf("#")+1));
      }
    },
    inviteProcess(str) {
      if(str.length != 11) return false;
      try{
        str = b64ToHex(str + "=");
      }catch(e){
        return false;
      }
      if(str.length != 16) return false;

      this.p2pID = str;
      this.p2pInviteLink = hexToB64(this.p2pID).substring(0,11);

      this.checkP2PGame();
      return true
    },
    hashProcess(){
      let h = window.location.hash.substring(1);
      this.inviteProcess(h);
    },
    async clearScores() {
      if(this.p2pMode) {
        this.p2pScores = [];
        return;
      }
      if(!this.globalResults) {
        this.localScores = [];
        return;
      }
      try {
        await axios.post(this.backendServer + 'clear_scores');
        await this.fetchScores();
      } catch (error) {
        console.error(error);
      }
    },
    populateWeapons() {
      for(let weapon of this.weapons) {
        this.weaponDict[weapon.id] = weapon;
      }
    },
    toggleResults() {
      this.globalResults = !this.globalResults;
    },
    async fetchScores() {
      if(!this.globalResults) return;
      try {
        const response = await axios.get(this.backendServer + 'get_scores');
        this.scores = response.data;
      } catch (error) {
        console.error(error);
      }
    },
    async fetchWeapons() {
      try {
        const response = await axios.get(this.backendServer + 'choices');
        this.weapons = response.data;
        this.populateWeapons();
      } catch (error) {
        console.error(error);
      }
    },
    async makeChoice(id) {
      if(this.p2pMode) {
        this.p2pSocket.send(JSON.stringify({"choice":id}));
        return;
      }
      try {
        this.computerChoiceId = 0;
        this.computerChoice = '';
        const response = await axios.post(this.backendServer + 'play', { player: id });
        this.showResult(
          response.data.results,
          this.weaponDict[response.data.computer],
          this.weaponDict[response.data.player]
        );
      } catch (error) {
        console.error(error);
      }
    },
    showResult(result, playerChoice, computerChoice) {
      this.result = result;
      this.computerChoice = playerChoice.name;
      this.yourChoice = computerChoice.name;
      this.computerChoiceId = playerChoice.id;
      this.yourChoiceId = computerChoice.id;
      this.resultRepresentation = this.getResultRepresentation(playerChoice.name, computerChoice.name);
      this.isShowResult = true;
      this.localScores.unshift(result);
      if(this.localScores.length>10) this.localScores = this.localScores.slice(0, 10);
      this.fetchScores();
    },
    closeResultModal() {
      this.isShowResult = false;
    },
    openChangeBackendModal() {
      this.isShowChangeBackendServer = true;
    },
    closeChangeBackendServerModal() {
      this.isShowChangeBackendServer = false;
    },
    async changeBackendServer() {
      try {
        axios.defaults.baseURL = this.backendServer;
        await this.fetchWeapons();
        this.isShowChangeBackendServer = false;
        this.cancelP2P();
      } catch (error) {
        console.error(error);
      }
    }
  }
}
</script>
<style>
body {
  height: 100vh;
  width: 100%;
  padding:0;
  margin:0;
  background: linear-gradient(-190deg,rgb(22, 23, 45),rgb(22, 23, 45) 60%,rgb(4, 6, 26)) #15172c;
  color: #acaec5;
  font-family:Muller,Ubuntu,"Droid Sans","Helvetica Neue",Arial,sans-serif,"Apple Color Emoji","Segoe UI Emoji","Segoe UI Symbol";
}
input {
  font-family:Muller,Ubuntu,"Droid Sans","Helvetica Neue",Arial,sans-serif,"Apple Color Emoji","Segoe UI Emoji","Segoe UI Symbol";
}
.modal{
  background: linear-gradient(-25deg,rgba(107,36,101,.54),rgba(107,36,101,.2) 40%,rgba(107,36,101,0)) #222040;
  border-radius: .6em;
  max-height:98vh;
  max-width:98vw;
}

.player {
  text-align: center;
}
.player>h1 {
  font-size: 1.6em;
  margin:0;
  height: 1.3em;
}
.player .weapon {
  margin: 1em auto;
  color: transparent;
  width:20vw;
  height: 20vw;
  background-size: cover;
}
.player .player-you-mark,
.player .waiting-player {
  opacity: .4;
  font-style: italic;
}
.backend-server{
  min-width: 32em;
}
.modal>h1{
  font-size: 1.3em;
  padding: .2em .6em;
  margin: 0 0 .3em 0;
  background-color: rgba(0,0,0,0.1);
  border-bottom:1px solid rgba(255,255,255,.05);
}
.modal>.content{
  position: relative;
  padding: .2em .8em;
}
.modal input{
  font-size: 1.1em;
  display: block;
  width: 100%;
  background: none;
  border: 0 none transparent;
  color: #fff;
  margin: .1em 0;
  padding: .2em 0;
  border-bottom: 2px solid rgba(148, 45, 139, 0.54);
}
.modal input:hover{
  outline:none;
  border-bottom: 2px solid rgba(221, 61, 208, 0.54);
}
.modal input:focus{
  outline:none;
  border-bottom: 2px solid rgba(223, 66, 145, 0.788);
}
.modal input:invalid{
  color:#f88;
  border-bottom: 2px solid rgba(255, 0, 0, 1);
}
.modal .content button {
  background: rgba(80, 150, 255, 0.1);
  font-size: 1em;
  margin: 0;
  padding: .1em .7em;
  border: 0 none transparent;
  cursor: pointer;
  color: #acaec5;
}
.modal .content button:hover {
  color: #fff;
}
.modal .footer{
  display: flex;
  justify-content: flex-end;
  padding: 0 .8em;
}
.modal .footer button{
  background: none;
  margin: .1em .2em;
  border: 0 none transparent;
  cursor: pointer;
  color: #acaec5;
  border-top: 2px solid rgba(148, 45, 139, 0.54);
}
.modal .footer button:hover{
  border-top: 2px solid rgba(221, 61, 208, 0.54);
  color: #fff;
}
.result .content{
  text-align: center;
}
.result img {
  max-height: 80vh;
}
.weapons .weapon {
  --translator-radial: -75%;
  --translator-base-x: -50%;
  --translator-base-y: -50%;
  --rotor-1: 45deg;
  --rotor-1-back:calc(var(--rotor-1)*-1);
  color:transparent;
  background-color: transparent;
  border: 0 none transparent;
  width: var(--weapon-size);
  height: var(--weapon-size);
  position: absolute;
  top: 50%;
  left: 50%;
  background-size: cover;
}
.weapons .weapon:hover{
  filter: contrast(170%);
}
.weapons>h1{
  font-size: 1.6em;
  margin:0;
  height: 1.3em;
}
.weapons-choices {
  --weapon-size: 7em;
  position: relative;
  height: calc(100% - (1.6em * 1.3));
  width: 100%;
}
.weapons {
  text-align: center;
  position: relative;
  height: 100%;
}
nav ul {
  display: flex;
  list-style: none;
  margin: 0;
  height: 100%;
  padding: 0;
}
nav li {
  margin: 0;
  padding: 0;
  flex: 1;
  text-align: center;
}
nav a {
  display: block;
  height: 100%;
  line-height: 1.6em;
  font-size: 1em;
  color:#acaec5;
  text-decoration: none;
  cursor: pointer;
}
nav a:hover{
  background-color: #20233e;
  color: #fff;
}
nav{
  display: block;
  position: relative;
  height: 1.8em;
}
main {
  overflow-y: auto;
  width: 100%;
  height: calc(100vh - 1.8em);
}
.game{
  display: grid;
  grid-template-columns: 1fr 2fr 1fr;
  height: calc(100vh - 6em);
}
.board{
  display: grid;
  grid-template: 1fr 1fr / 1fr 2em;
  margin: 0 auto;
  max-width: 40em;
}
.board>h1{
  grid-column: 1/3;
  text-align: center;
  font-size: 1.5em;
  margin: 0;
  padding: 0;
}
.board>button{
  width: 1.6em;
  font-size: 1.1em;
  color: #acaec5;
  cursor:pointer;
  background:none;
  border: 0 none transparent;
}
.board>button:hover{
  color: #fff;
}
.scores{
  overflow-x: hidden;
  white-space: nowrap;
  font-size: 1.5em;
}
.score {
  text-align: center;
  display: inline-block;
  margin: 0 .1em;
  padding: 0 .3em;
}
.score:nth-child(7){
  opacity: .8;
}
.score:nth-child(8){
  opacity: .5;
}
.score:nth-child(9){
  opacity: .3;
}
.score:nth-child(10){
  opacity: .05;
}





.weapons .weapon-1 {
  transform:translate(var(--translator-base-x), var(--translator-base-y)) rotate(var(--rotor-1)) rotate(0deg) translate(var(--translator-radial), var(--translator-radial))  rotate(-0deg) rotate(var(--rotor-1-back));
}
.weapons .weapon-2 {
  transform:translate(var(--translator-base-x), var(--translator-base-y)) rotate(var(--rotor-1)) rotate(72deg) translate(var(--translator-radial), var(--translator-radial)) rotate(-72deg) rotate(var(--rotor-1-back));
}
.weapons .weapon-3 {
  transform:translate(var(--translator-base-x), var(--translator-base-y)) rotate(var(--rotor-1)) rotate(144deg) translate(var(--translator-radial), var(--translator-radial)) rotate(-144deg) rotate(var(--rotor-1-back));
}
.weapons .weapon-4 {
  transform:translate(var(--translator-base-x), var(--translator-base-y)) rotate(var(--rotor-1)) rotate(216deg) translate(var(--translator-radial), var(--translator-radial)) rotate(-216deg) rotate(var(--rotor-1-back));
}
.weapons .weapon-5 {
  transform:translate(var(--translator-base-x), var(--translator-base-y)) rotate(var(--rotor-1)) rotate(288deg) translate(var(--translator-radial), var(--translator-radial)) rotate(-288deg) rotate(var(--rotor-1-back));
}
.weapon-1 {
  background-image: url(./images/rock.png);
}
.weapon-2 {
  background-image: url(./images/paper.png);
}
.weapon-3 {
  background-image: url(./images/scissors.png);
}
.weapon-4 {
  background-image: url(./images/lizard.png);
}
.weapon-5 {
  background-image: url(./images/spock.png);
}
</style>
