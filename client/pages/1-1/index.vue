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
      camera: {},
    }
  },
  created() {
    socket.emit('setup')
  },
  mounted() {
    this.$refs.bg.width = 16 * 16
    this.$refs.bg.height = 16 * 16
    this.screen.context = this.$refs.bg.getContext('2d')
    this.background = this.getCanvas(256, 256)

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
          window.camera = data.camera.position
          this.camera = data.camera
          this.tile = this.getCanvas(256 + this.camera.position.x, 256)
          e.ranges.forEach((range) => {
            this.drawBG(block, range, e.isCollide)
          })

          this.background
            .getContext('2d')
            .drawImage(
              this.tile,
              -data.camera.position.x,
              -data.camera.position.y
            )
        })
      })
    })

    this.loadImage(require('@/assets/images/characters.gif')).then((image) => {
      socket.on('drawMario', (mario) => {
        const buffer = this.draw(
          image,
          mario.action.image.x,
          mario.action.image.y,
          mario.action.size.width,
          mario.action.size.height,
          mario.movement.includes('left')
        )

        this.background
          .getContext('2d')
          .drawImage(
            buffer,
            mario.position.x - this.camera.position.x,
            mario.position.y - this.camera.position.y
          )
        this.background.getContext('2d').strokeStyle = 'blue'
        this.background.getContext('2d').beginPath()
        this.background
          .getContext('2d')
          .rect(
            mario.position.x - this.camera.position.x,
            mario.position.y - this.camera.position.y,
            this.tileSize,
            this.tileSize
          )
        this.background.getContext('2d').stroke()
        this.screen.context.drawImage(this.background, 0, 0)
        socket.emit('fall')
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
          socket.emit('jump')
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
    draw(image, x, y, width, height, isFlip = false) {
      const buffer = this.getCanvas(width, height)
      const context = buffer.getContext('2d')

      if (isFlip) {
        context.scale(-1, 1)
        context.translate(-width, 0)
      }

      context.drawImage(image, x, y, width, height, 0, 0, width, height)
      return buffer
    },
    drawBG(block, range, isCollide) {
      for (let i = range.x1; i < range.x2; i++) {
        for (let j = range.y1; j < range.y2; j++) {
          this.tile.getContext('2d').drawImage(block, i * 16, j * 16)
          if (isCollide) {
            this.tile.getContext('2d').strokeStyle = 'red'
            this.tile.getContext('2d').beginPath()
            this.tile
              .getContext('2d')
              .rect(i * 16, j * 16, this.tileSize, this.tileSize)
            this.tile.getContext('2d').stroke()
          }
        }
      }
    },
  },
}
</script>
