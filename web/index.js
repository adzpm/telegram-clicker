let Application = Vue.createApp({
    data: function () {
        return {
            telegram_data: {},

            tmp_btn_count: [1],
            user_data: {
                coins: 0,
            },
        }
    },

    methods: {
        Enter(telegram_id) {
            let url = this.CurrentURL + '/enter' + '?telegram_id=' + telegram_id

            axios.get(url).then(response => {
                console.log(response.data)
                this.user_data = response.data
            }).catch(error => {
                console.log(error)
            })
        },

        Click(telegram_id, product_id) {
            let url = this.CurrentURL + '/click' + '?telegram_id=' + telegram_id + '&product_id=' + product_id

            axios.get(url).then(response => {
                console.log(response.data)
                this.user_data = response.data
            }).catch(error => {
                console.log(error)
            })
        },

        BuyProduct(telegram_id, product_id) {
            let url = this.CurrentURL + '/buy' + '?telegram_id=' + telegram_id + '&product_id=' + product_id

            axios.get(url).then(response => {
                console.log(response.data)
                this.user_data = response.data
            }).catch(error => {
                console.log(error)
            })
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
            setTimeout(() => {}, 10)
        }

        this.telegram_data = {...window.Telegram?.WebApp?.initDataUnsafe}

        this.Enter(this.TelegramID)
    },
})

const vm = Application.mount('#app')