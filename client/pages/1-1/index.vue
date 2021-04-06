<template>
  <div>
    <canvas ref="bg" style="background-color: white" />
  </div>
</template>

<script>
import socket from '@/plugins/socket.io'

export default {
  components: {},
  data() {
    return {
      screen: {
        context: {},
      },
      background: {},
      width: 16,
      height: 16,
    }
  },
  created() {
    socket.emit('setup')
  },
  mounted() {
    this.$refs.bg.width = window.innerWidth
    this.$refs.bg.height = window.innerHeight
    this.screen.context = this.$refs.bg.getContext('2d')
    this.background = this.getCanvas(window.innerWidth, window.innerHeight)

    socket.on('setup', (data) => {
      const tileSize = 16
      this.loadImage(require('@/assets/images/bg_tiles.png')).then((image) => {
        const block = this.draw(
          image,
          data.position.x * tileSize,
          data.position.y * tileSize,
          tileSize,
          tileSize,
          0,
          0,
          tileSize,
          tileSize
        )
        for (const range of data.ranges) {
          this.drawBG(this.background.getContext('2d'), block, range)
          this.screen.context.drawImage(this.background, 0, 0)
        }
      })
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
    getCanvas(width, height) {
      const buffer = document.createElement('canvas')
      buffer.width = width
      buffer.height = height
      return buffer
    },
    draw(image, x, y, width, height) {
      const buffer = this.getCanvas(width, height)
      buffer
        .getContext('2d')
        .drawImage(image, x, y, width, height, 0, 0, width, height)
      return buffer
    },
    drawBG(bg, block, boundary) {
      for (let i = boundary.x1; i < boundary.x2; i++) {
        for (let j = boundary.y1; j < boundary.y2; j++) {
          bg.drawImage(block, i * 16, j * 16)
        }
      }
    },
    buildMario(mario, position) {
      this.loadImage(require('@/assets/images/characters.gif')).then(
        (image) => {
          const buffer = this.draw(
            image,
            mario.x,
            mario.y,
            mario.width,
            mario.height
          )
          this.background
            .getContext('2d')
            .drawImage(
              buffer,
              position.x + mario.width,
              position.y + mario.height
            )
          this.screen.context.drawImage(this.background, 0, 0)
        }
      )
    },
  },
}
</script>
