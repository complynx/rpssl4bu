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
        this.result = response.data.results;
        this.computerChoice = this.weaponDict[response.data.computer].name;
        this.yourChoice = this.weaponDict[response.data.player].name;
        this.isShowResult = true;
        await this.fetchScores();
      } catch (error) {
        console.error(error);
      }
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
