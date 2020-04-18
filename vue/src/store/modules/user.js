import { logout } from '@/api/user'
import { login, getInfo } from '@/api/weagent'
import { getToken, setToken, removeToken } from '@/utils/auth'
import router, { resetRouter } from '@/router'

const state = {
    token: getToken(),
    name: '',
    avatar: '',
    introduction: '',
    roles: []
}

const mutations = {
    SET_TOKEN: (state, token) => {
        state.token = token
    },
    SET_INTRODUCTION: (state, introduction) => {
        state.introduction = introduction
    },
    SET_NAME: (state, name) => {
        state.name = name
    },
    SET_AVATAR: (state, avatar) => {
        state.avatar = avatar
    },
    SET_ROLES: (state, roles) => {
        state.roles = roles
    }
}

const actions = {
    // user login
    login ({ commit }, userInfo) {
        const { username, password } = userInfo
        return new Promise((resolve, reject) => {
            var params = {
                account: username,
                password: password
            }
            var data = JSON.stringify(params)
            login(data).then(response => {
                commit('SET_TOKEN', response.token)
                setToken(response.token)
                resolve()
            }).catch(error => {
                reject(error)
            })
        })
    },

    // get user info
    getInfo ({ commit, state }) {
        return new Promise((resolve, reject) => {
            getInfo(state.token).then(response => {
                if (!response) {
                    reject('Verification failed, please Login again.')
                }

                var data = {
                    roles: response.role,
                    name: response.nickname,
                    avatar: response.avatarurl,
                    introduction: "yeah yeah",
                }

                // roles must be a non-empty array
                if (!data.roles || data.roles.length <= 0) {
                    reject('getInfo: roles must be a non-null array!')
                }

                commit('SET_ROLES', data.roles)
                commit('SET_NAME', data.name)
                commit('SET_AVATAR', data.avatar)
                commit('SET_INTRODUCTION', data.introduction)
                resolve(data)
            }).catch(error => {
                reject(error)
            })
        })
    },

    // remove token
    resetToken ({ commit }) {
        return new Promise(resolve => {
            commit('SET_TOKEN', '')
            commit('SET_ROLES', [])
            removeToken()
            resolve()
        })
    },

    // user logout
    logout ({ commit, state, dispatch }) {
        return new Promise((resolve, reject) => {
            logout(state.token).then(() => {
                commit('SET_TOKEN', '')
                commit('SET_ROLES', [])
                removeToken()
                resetRouter()

                // reset visited views and cached views
                // to fixed https://github.com/PanJiaChen/vue-element-admin/issues/2485
                dispatch('tagsView/delAllViews', null, { root: true })

                resolve()
            }).catch(error => {
                reject(error)
            })
        })
    }
}

export default {
    namespaced: true,
    state,
    mutations,
    actions
}
