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
      tileSize: 16,
    }
  },
  created() {
    socket.emit('setup')
    socket.emit('draw')
  },
  mounted() {
    this.$refs.bg.width = window.innerWidth
    this.$refs.bg.height = window.innerHeight
    this.screen.context = this.$refs.bg.getContext('2d')
    this.background = this.getCanvas(window.innerWidth, window.innerHeight)

    this.loadImage(require('@/assets/images/bg_tiles.png')).then((image) => {
      socket.on('draw', (data) => {
        data.backgrounds.forEach((e) => {
          const block = this.draw(
            image,
            e.position.x * this.tileSize,
            e.position.y * this.tileSize,
            this.tileSize,
            this.tileSize,
            0,
            0,
            this.tileSize,
            this.tileSize
          )
          e.ranges.forEach((range) => {
            this.drawBG(block, range)
          })
        })

        data.colliders.forEach((collider) => {
          this.background.getContext('2d').strokeStyle = 'red'
          if (collider.name === 'ground') {
            this.background.getContext('2d').beginPath()
            this.background
              .getContext('2d')
              .rect(
                collider.tile.x,
                collider.tile.y,
                this.tileSize,
                this.tileSize
              )
            this.background.getContext('2d').stroke()
          }
        })
      })
    })

    this.loadImage(require('@/assets/images/characters.gif')).then((image) => {
      socket.on('drawMario', (mario) => {
        const buffer = this.draw(
          image,
          mario.x,
          mario.y,
          mario.width,
          mario.height
        )
        this.background
          .getContext('2d')
          .drawImage(buffer, mario.position.x, mario.position.y)
        this.background.getContext('2d').strokeStyle = 'blue'
        this.background.getContext('2d').beginPath()
        this.background
          .getContext('2d')
          .rect(
            mario.position.x,
            mario.position.y,
            this.tileSize,
            this.tileSize
          )
        this.background.getContext('2d').stroke()
        this.screen.context.drawImage(this.background, 0, 0)
        socket.emit('draw')
      })
    })

    window.addEventListener('keydown', function (e) {
      switch (e.code) {
        case 'KeyD':
          socket.emit('right')
          break
        case 'KeyA':
          socket.emit('left')
          break
        case 'KeyW':
          socket.emit('up')
          break
      }
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
    drawBG(block, boundary) {
      for (let i = boundary.x1; i < boundary.x2; i++) {
        for (let j = boundary.y1; j < boundary.y2; j++) {
          this.background.getContext('2d').drawImage(block, i * 16, j * 16)
        }
      }
    },
  },
}
</script>
