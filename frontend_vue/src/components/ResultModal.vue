<template>
  <div class="result-modal" v-if="isVisible">
    <div class="result-modal-content">
        <template v-if="result === 'win'">
            <h1>You win!</h1>
        </template>
        <template v-if="result === 'lose'">
            <h1>You've lost</h1>
        </template>
        <template v-if="result === 'tie'">
            <h1>It's a tie</h1>
        </template>
    </div>
  </div>
</template>


<script>
export default {
  name: 'ResultModal',
  props: {
    isVisible: {
      type: Boolean,
      default: false
    },
    result: {
      type: String
    }
  }
}
</script>

<style>
.result-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 999;
}

.result-modal-content {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background-color: white;
  padding: 20px;
  border-radius: 10px;
}
</style>
