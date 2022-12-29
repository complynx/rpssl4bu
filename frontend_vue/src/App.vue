<template>
  <nav>
    <ul>
      <li>
        <a href="#">Play P2P</a>
      </li>
      <li>
        <a @click="openChangeBackendModal">Change backend server</a>
      </li>
    </ul>
  </nav>
  <main>
    <div class="game">
      <div class="player">
        <!-- Add content for the player section here -->
      </div>
      <div class="weapons">
        <h1>Choose your weapon</h1>
        <div class="weapons-choices">
          <button class="weapon" v-for="weapon in weapons" :key="weapon.id" :class="'weapon-'+weapon.id" @click="makeChoice(weapon.id)">
            {{ weapon.name }}
          </button>
        </div>
      </div>
      <div class="computer">
        <!-- Add content for the computer section here -->
      </div>
    </div>
    <div class="board">
      <h1>Last scores:</h1>
      <div class="scores">
        <div class="score" v-for="(score, index) in scores" :key="index">
          {{ score }}
        </div>
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
      <template v-if="result === 'win'">
          <h1>You win!</h1>
          <div class="content">
            <div>Computer chose {{computerChoice}}</div>
            <div>You chose {{yourChoice}}</div>
            <template v-if="resultRepresentation !== null">
              <div>{{resultRepresentation.text}}</div>
              <img :src="'./images/'+resultRepresentation.img" />
            </template>
          </div>
      </template>
      <template v-if="result === 'lose'">
          <h1>You've lost</h1>
          <div class="content">
            <div>Computer chose {{computerChoice}}</div>
            <div>You chose {{yourChoice}}</div>
            <template v-if="resultRepresentation !== null">
              <div>{{resultRepresentation.text}}</div>
              <img :src="'./images/'+resultRepresentation.img" />
            </template>
          </div>
      </template>
      <template v-if="result === 'tie'">
          <h1>It's a tie</h1>
          <div class="content">
            <div>Computer chose {{computerChoice}}</div>
            <div>You chose {{yourChoice}}</div>
            <template v-if="resultRepresentation !== null">
              <div>{{resultRepresentation.text}}</div>
              <img :src="'./images/'+resultRepresentation.img" />
            </template>
          </div>
      </template>
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
</template>

<script>
import axios from 'axios';

export default {
  name: 'App',
  components: {
  },
  data() {
    return {
      // backendServer: '/backend/',
      backendServer: 'http://localhost:80/backend/',
      // backendServer: 'https://codechallenge.boohma.com/',
      isShowChangeBackendServer: false,
      isShowResult: false,
      computerChoice: '',
      yourChoice: '',
      result: '',
      resultRepresentation: null,
      scores: [],
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
      }
      return null;
    },
    async clearScores() {
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
    async fetchScores() {
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
      try {
        const response = await axios.post(this.backendServer + 'play', { player: id });
        this.showResult(
          response.data.results,
          this.weaponDict[response.data.computer].name,
          this.weaponDict[response.data.player].name
        );
      } catch (error) {
        console.error(error);
      }
    },
    showResult(result, playerChoice, computerChoice) {
      this.result = result;
      this.computerChoice = playerChoice;
      this.yourChoice = computerChoice;
      this.resultRepresentation = this.getResultRepresentation(playerChoice, computerChoice);
      console.log(this.resultRepresentation);
      this.isShowResult = true;
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
.weapon {
  --translator-radial: -75%;
  --translator-base-x: -50%;
  --translator-base-y: -50%;
  --rotor-1: 45deg;
  --rotor-1-back:calc(var(--rotor-1)*-1);
  color:transparent;
  background: none;
  border: 0 none transparent;
  width: var(--weapon-size);
  height: var(--weapon-size);
  position: absolute;
  top: 50%;
  left: 50%;
  background-size: cover;
}
.weapon:hover{
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





.weapon-1 {
  transform:translate(var(--translator-base-x), var(--translator-base-y)) rotate(var(--rotor-1)) rotate(0deg) translate(var(--translator-radial), var(--translator-radial))  rotate(-0deg) rotate(var(--rotor-1-back));
  background-image: url(./images/rock.png);
}
.weapon-2 {
  transform:translate(var(--translator-base-x), var(--translator-base-y)) rotate(var(--rotor-1)) rotate(72deg) translate(var(--translator-radial), var(--translator-radial)) rotate(-72deg) rotate(var(--rotor-1-back));
  background-image: url(./images/paper.png);
}
.weapon-3 {
  transform:translate(var(--translator-base-x), var(--translator-base-y)) rotate(var(--rotor-1)) rotate(144deg) translate(var(--translator-radial), var(--translator-radial)) rotate(-144deg) rotate(var(--rotor-1-back));
  background-image: url(./images/scissors.png);
}
.weapon-4 {
  transform:translate(var(--translator-base-x), var(--translator-base-y)) rotate(var(--rotor-1)) rotate(216deg) translate(var(--translator-radial), var(--translator-radial)) rotate(-216deg) rotate(var(--rotor-1-back));
  background-image: url(./images/lizard.png);
}
.weapon-5 {
  transform:translate(var(--translator-base-x), var(--translator-base-y)) rotate(var(--rotor-1)) rotate(288deg) translate(var(--translator-radial), var(--translator-radial)) rotate(-288deg) rotate(var(--rotor-1-back));
  background-image: url(./images/spock.png);
}
</style>
