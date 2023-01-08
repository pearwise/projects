<script setup lang="ts">
// import { ref } from 'vue'

// defineProps<{ msg: string }>()

// const count = ref(0)
</script>

<template>
  <div class="box">
    <div class="top-box">
      <img class="bg-img" src="../assets/xinguan.jpg" alt="" />
      <div class="title-text">
        <h1>科学防护 共度难关</h1>
        <h2>肺炎疫情实时动态播报</h2>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.bg-img {
  width:100%;
}
.top-box {
  position: relative;
  .bg-img{
    position: absolute;
  }
  .title-text {
    z-index: 2;
    color: #fff;
  }
}
</style>
