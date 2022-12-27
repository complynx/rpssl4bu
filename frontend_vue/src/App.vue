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
        <button class="weapon" v-for="weapon in weapons" :key="weapon.id" :class="'weapon-'+weapon.id" @click="makeChoice(weapon.id)">
          {{ weapon.name }}
        </button>
      </div>
      <div class="computer">
        <!-- Add content for the computer section here -->
      </div>
    </div>
    <div class="board">
      <h1>Last scores:</h1>
      <div class="score" v-for="(score, index) in scores" :key="index">
        {{ score }}
      </div>
      
      <button @click="clearScores()">
        clear
      </button>
    </div>
  </main>
  <Modal
    v-model="isShowResult"
    :close="closeResultModal"
  >
    <div class="modal">
      <template v-if="result === 'win'">
          <h1>You win!</h1>
          <div>Computer chose {{computerChoice}}</div>
          <div>You chose {{yourChoice}}</div>
      </template>
      <template v-if="result === 'lose'">
          <h1>You've lost</h1>
          <div>Computer chose {{computerChoice}}</div>
          <div>You chose {{yourChoice}}</div>
      </template>
      <template v-if="result === 'tie'">
          <h1>It's a tie</h1>
          <div>Computer chose {{computerChoice}}</div>
          <div>You chose {{yourChoice}}</div>
      </template>
    </div>
  </Modal>
  <Modal
    v-model="isShowChangeBackendServer"
    :close="closeChangeBackendServerModal"
  >
    <div class="modal">
      <h1>Change backend server</h1>
      <p>Enter the new backend server address:</p>
      <input v-model="backendServer" type="text">
      <button @click="changeBackendServer">Save</button>
      <button @click="closeChangeBackendServerModal">
        Cancel
      </button>
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
  },
  methods: {
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
.weapon {
  color:transparent;
  width:5em;
  height:5em;
  position: absolute;
  top: 50%;
  left: 50%;
  --translator: -80%;
  --rotor-1: 45deg;
  --rotor-1-back:calc(var(--rotor-1)*-1);
  background-size: cover;
}
.weapon-1 {
  transform: rotate(var(--rotor-1)) rotate(0deg) translate(var(--translator), var(--translator))  rotate(-0deg) rotate(var(--rotor-1-back));
  background-image: url('./images/rock.png');
}
.weapon-2 {
  transform: rotate(var(--rotor-1)) rotate(72deg) translate(var(--translator), var(--translator)) rotate(-72deg) rotate(var(--rotor-1-back));
  background-image: url('./images/paper.png');
}
.weapon-3 {
  transform: rotate(var(--rotor-1)) rotate(144deg) translate(var(--translator), var(--translator)) rotate(-144deg) rotate(var(--rotor-1-back));
  background-image: url('./images/scissors.png');
}
.weapon-4 {
  transform: rotate(var(--rotor-1)) rotate(216deg) translate(var(--translator), var(--translator)) rotate(-216deg) rotate(var(--rotor-1-back));
  background-image: url('./images/lizard.png');
}
.weapon-5 {
  transform: rotate(var(--rotor-1)) rotate(288deg) translate(var(--translator), var(--translator)) rotate(-288deg) rotate(var(--rotor-1-back));
  background-image: url('./images/spock.png');
}
.weapons {
  text-align: center;
  position: relative;
  height: 10em;
}

.game {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
}

nav ul {
  display: flex;
  list-style: none;
  margin: 0;
  padding: 0;
}
nav li {
  flex: 1;
  text-align: center;
}
nav a {
  display: block;
  padding: 10px;
  font-size: 18px;
  color:rgb(37, 37, 37);
  text-decoration: none;
  cursor: pointer;
}
nav a:hover{
  background-color: rgb(224, 224, 224);
}

#app {
  font-family: Helvetica, Arial, sans-serif;
}
</style>
