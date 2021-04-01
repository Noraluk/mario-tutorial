<template>
  <div>
    <canvas ref="bg" :width="background.width" :height="background.height" />
  </div>
</template>

<script>
import socket from '../plugins/socket.io'

export default {
  components: {},
  data() {
    return {
      background: {
        context: {},
        width: 200,
        height: 200,
      },
      width: 16,
      height: 16,
    }
  },
  created() {
    socket.emit('start')
  },
  mounted() {
    socket.on('onStart', (data) => {
      this.buildScene(this.$refs.bg.getContext('2d'), data.tile, data.boundary)
    })
  },
  methods: {
    loadImage(url) {
      return new Promise((resolve) => {
        const image = new Image()
        image.addEventListener('load', () => {
          resolve(image)
        })
        image.src = url
      })
    },
    buildScene(context, tile, boundary) {
      this.loadImage(require('../assets/images/bg_tiles.png')).then((image) => {
        const buffer = document.createElement('canvas')
        buffer.height = tile.height
        buffer.width = tile.width
        buffer
          .getContext('2d')
          .drawImage(
            image,
            tile.width * tile.x,
            tile.height * tile.y,
            tile.width,
            tile.height,
            0,
            0,
            tile.width,
            tile.height
          )
        for (let i = boundary.x1; i < boundary.x2; i++) {
          for (let j = boundary.y1; j < boundary.y2; j++) {
            context.drawImage(buffer, i * tile.width, j * tile.height)
          }
        }
      })
    },
  },
}
</script>
