/* eslint-disable-next-line */
import Vue from 'vue'
import businessSelector from './selector/business.vue'
import clipboardSelector from './selector/clipboard.vue'
import selector from './selector/selector.vue'
import details from './details/details.vue'
import form from './form/form.vue'
import formMultiple from './form/form-multiple.vue'
import bool from './form/bool.vue'
import boolInput from './form/bool-input.vue'
import date from './form/date.vue'
import dateRange from './form/date-range.vue'
import time from './form/time.vue'
import int from './form/int.vue'
import float from './form/float.vue'
import longchar from './form/longchar.vue'
import singlechar from './form/singlechar.vue'
import timezone from './form/timezone.vue'
import enumeration from './form/enum.vue'
import objuser from './form/objuser.vue'
import tree from './tree/tree.vue'
import resize from './other/resize.vue'
import collapseTransition from './transition/collapse.js'
import collapse from './collapse/collapse'
import dotMenu from './dot-menu/dot-menu.vue'
import input from './form/input.vue'
import searchInput from './form/search-input.vue'
import inputSelect from './selector/input-select.vue'
const install = (Vue, opts = {}) => {
    const components = [
        businessSelector,
        clipboardSelector,
        selector,
        details,
        form,
        formMultiple,
        bool,
        boolInput,
        date,
        dateRange,
        time,
        int,
        float,
        longchar,
        singlechar,
        timezone,
        enumeration,
        objuser,
        tree,
        resize,
        collapseTransition,
        collapse,
        dotMenu,
        input,
        searchInput,
        inputSelect
    ]
    components.forEach(component => {
        Vue.component(component.name, component)
    })
}

export default {
    install,
    businessSelector,
    clipboardSelector,
    selector,
    details,
    form,
    formMultiple,
    bool,
    boolInput,
    date,
    dateRange,
    time,
    int,
    float,
    longchar,
    singlechar,
    timezone,
    enumeration,
    objuser,
    tree,
    resize,
    collapseTransition,
    dotMenu,
    input,
    searchInput,
    inputSelect
}
