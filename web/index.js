let Application = Vue.createApp({
    data: function () {
        return {
            started_from_telegram: false,
            telegram_data: null,
            game_data: null,
            error: null,
        }
    },

    methods: {
        Enter(telegram_id) {
            let url = this.CurrentURL + '/enter' + '?telegram_id=' + telegram_id

            axios.get(url).then(response => {
                console.log(response.data)
                this.game_data = response.data
            }).catch(error => {
                this.ShowError(error.response.data)
            })
        },

        Click(telegram_id, product_id) {
            let url = this.CurrentURL + '/click' + '?telegram_id=' + telegram_id + '&product_id=' + product_id

            axios.get(url).then(response => {
                console.log(response.data)
                this.game_data = response.data
            }).catch(error => {
                this.ShowError(error.response.data)
            })
        },

        BuyProduct(telegram_id, product_id) {
            let url = this.CurrentURL + '/buy' + '?telegram_id=' + telegram_id + '&product_id=' + product_id

            axios.get(url).then(response => {
                console.log(response.data)
                this.game_data = response.data
            }).catch(error => {
                this.ShowError(error.response.data)
            })
        },

        // non-api methods

        ShowError(err) {
            this.error = err
            setTimeout(() => {
                this.error = null
            }, 2500)
        },

        PopEffect: function (e) {
            for (let i = 0; i < 5; i++) {
                this.createParticle(e.clientX, e.clientY);
            }
        },

        createParticle: function (x, y) {
            let destinationX = x + (Math.random() - 0.5) * 2 * 75,
                destinationY = y + (Math.random() - 0.5) * 2 * 75,
                particle = document.createElement('particle'),
                symbols = ['$', '↑', '%']

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

        FormatNumber(num) {
            const units = [
                {value: 1e36, suffix: ' INFINITY'},
                {value: 1e33, suffix: ' QUINTILLIARD'},
                {value: 1e30, suffix: ' QUINTILLION'},
                {value: 1e27, suffix: ' QUADRILLIARD'},
                {value: 1e24, suffix: ' QUADRILLION'},
                {value: 1e21, suffix: ' TRILLIARD'},
                {value: 1e18, suffix: ' TRILLION'},
                {value: 1e15, suffix: ' BILLIARD'},
                {value: 1e12, suffix: ' BILLION'},
                {value: 1e9, suffix: ' MILLIARD'},
                {value: 1e6, suffix: ' MILLION'},
            ];

            for (let unit of units) {
                if (num >= unit.value) {
                    return (num / unit.value).toFixed(3).replace(/\.?0+$/, '') + unit.suffix;
                }
            }

            return num.toString();
        }
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