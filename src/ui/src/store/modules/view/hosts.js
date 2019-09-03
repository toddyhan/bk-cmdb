const defaultParams = {
    ip: {
        flag: 'bk_host_innerip|bk_host_outer',
        exact: 0,
        data: []
    }
}
const state = {
    filterList: [],
    filterIP: null,
    filterParams: {
        ...defaultParams
    },
    collection: null,
    collectionList: []
}

const getters = {
    isCollection: state => !!state.collection
}

const mutations = {
    setFilterList (state, list) {
        state.filterList = list
    },
    setFilterIP (state, IP) {
        state.filterIP = IP
    },
    setFilterParams (state, params) {
        state.filterParams = params
    },
    setCollectionList (state, list) {
        state.collectionList = list
    },
    setCollection (state, collection) {
        state.collection = collection
    },
    addCollection (state, collection) {
        state.collectionList.push(collection)
    },
    updateCollection (state, updatedData) {
        Object.assign(state.collection, updatedData)
    },
    deleteCollection (state, id) {
        state.collectionList = state.collectionList.filter(collection => collection.id === id)
    },
    clearFilter (state) {
        state.filterList = []
        state.filterIP = null
        state.filterParams = {
            ...defaultParams
        }
        state.collection = null
    }
}

export default {
    namespaced: true,
    state,
    getters,
    mutations
}
