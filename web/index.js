let Application = Vue.createApp({
    data: function () {
        return {
            telegram_data: {},

            game_data: {
                "user_id": 4,
                "telegram_id": 9876543210,
                "last_seen": 1719341593,
                "current_coins": 240,
                "current_gold": 0,
                "products": {
                    "1": {
                        "id": 1,
                        "name": "BTC Exchange",
                        "image_url": "https://www.svgrepo.com/show/289569/exchange-euro.svg",
                        "upgrade_price": 4633,
                        "coins_per_click": 13,
                        "coins_per_minute": 13,
                        "current_level": 1,
                        "max_level": 100
                    },
                    "2": {
                        "id": 2,
                        "name": "Currency mapping",
                        "image_url": "https://www.svgrepo.com/show/430187/solution-bulb-concept.svg",
                        "upgrade_price": 2155,
                        "coins_per_click": 0,
                        "coins_per_minute": 0,
                        "current_level": 0,
                        "max_level": 100
                    },
                    "3": {
                        "id": 3,
                        "name": "Crypto Stock",
                        "image_url": "https://www.svgrepo.com/show/493504/stockfx-chart.svg",
                        "upgrade_price": 17630,
                        "coins_per_click": 0,
                        "coins_per_minute": 0,
                        "current_level": 0,
                        "max_level": 100
                    }
                }
            },
        }
    },

    methods: {
        Enter(telegram_id) {
            let url = this.CurrentURL + '/enter' + '?telegram_id=' + telegram_id

            axios.get(url).then(response => {
                console.log(response.data)
                this.game_data = response.data
            }).catch(error => {
                console.log(error)
            })
        },

        Click(telegram_id, product_id) {
            let url = this.CurrentURL + '/click' + '?telegram_id=' + telegram_id + '&product_id=' + product_id

            axios.get(url).then(response => {
                console.log(response.data)
                this.game_data = response.data
            }).catch(error => {
                console.log(error)
            })
        },

        BuyProduct(telegram_id, product_id) {
            let url = this.CurrentURL + '/buy' + '?telegram_id=' + telegram_id + '&product_id=' + product_id

            axios.get(url).then(response => {
                console.log(response.data)
                this.game_data = response.data
            }).catch(error => {
                console.log(error)
            })
        },

        // non-api methods

        PopEffect: function (e) {
            for (let i = 0; i < 5; i++) {
                this.createParticle(e.clientX, e.clientY);
            }
        },

        createParticle: function (x, y) {
            let destinationX = x + (Math.random() - 0.5) * 2 * 75,
                destinationY = y + (Math.random() - 0.5) * 2 * 75,
                particle = document.createElement('particle'),
                symbols = ['$', '$', '$']

            document.body.appendChild(particle);
            particle.innerHTML = symbols[Math.floor(Math.random() * symbols.length)]
            particle.style.fontSize = `${Math.random() * 24 + 10}px`;

            let animation = particle.animate([
                {transform: `translate(-50%, -50%) translate(${x}px, ${y}px)`, opacity: 1},
                {transform: `translate(${destinationX}px, ${destinationY}px)`, opacity: 0}
            ], {duration: Math.floor(Math.random() * 100 + 1000), easing: 'ease-out'})

            animation.onfinish = () => {
                particle.remove()
            }
        },
    },

    computed: {
        CurrentURL() {
            let url = window.location.href
            return url.substring(0, url.lastIndexOf('/'))
        },

        TelegramID() {
            return this.telegram_data?.user?.id ?? 9876543210
        },

        TelegramName() {
            let name = this.telegram_data?.user?.first_name

            if (this.telegram_data?.user?.last_name) {
                name += ' ' + this.telegram_data?.user?.last_name
            }

            return name ?? 'John Doe'
        }
    },

    mounted: function () {
        while (window.Telegram.WebApp.initDataUnsafe === undefined) {
            setTimeout(() => {
            }, 10)
        }

        this.telegram_data = {...window.Telegram?.WebApp?.initDataUnsafe}

        this.Enter(this.TelegramID)
    },
})

const vm = Application.mount('#app')