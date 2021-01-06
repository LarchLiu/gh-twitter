<template>
  <div>
    <div :id="id_name">
      <div :class="(!childClassName && fixed) ? styleClassName : ''" :style="(widthSelf && !childClassName && fixed) ? {width: fixed_element_width + 'px'} : ''">
        <slot />
      </div>
    </div>
    <div v-if="fixed" :style="{paddingBottom: fixed_element_height + 'px'}" />
  </div>
</template>

<script>
import _ from 'lodash'
/**
   * 支持传递组件的内部 class name 来 fixed 内部组件
   * 同一页面如有多个 FixedFooter 请为每个组件传递 idName, 因为是通过 id 找到具体的组件, 如果 id 相同多个组件无法区分
   * 单个 FixedHeader 可不用传递 idName, 只需在页面定义默认的名为 fixed-footer 的样式, 不考虑层级关系可以定义 /deep/.fixed-footer
   * 不用默认样式名称需传入自定义的 styleClassName
   */

export default {
  name: 'FixedFooter',
  props: {
    // 需要 fixed 组件 id 名称，用于定位组件
    idName: {
      type: String,
      default: 'element'
    },
    // 组件 fixed 时的样式名称
    styleClassName: {
      type: String,
      default: 'fixed-footer'
    },
    // 需要 fixed 组件中的子组件 class name, 用于定位到此子组件
    childClassName: {
      type: String,
      default: null
    },
    checkFixed: {
      type: Number,
      default: 0
    },
    widthSelf: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      id_name_pre: 'fixed-',
      fixed: false,
      fixed_element: null,
      fixed_element_width: 0,
      fixed_element_height: 0
    }
  },
  computed: {
    id_name () {
      return this.id_name_pre + this.idName
    }
  },
  watch: {
    fixed (val) {
      if (this.childClassName && this.fixed_element) {
        val ? this.addFixedClass(this.fixed_element, this.styleClassName) : this.removeFixedClass(this.fixed_element, this.styleClassName)
      }
      this.$emit('change', val)
    },
    checkFixed () {
      this.handleScroll()
    }
  },
  mounted () {
    window.addEventListener('scroll', _.throttle(this.handleScroll, 50))
    const element = document.querySelector(`#${this.id_name}`)
    if (this.childClassName) {
      const childElement = element.querySelector(`.${this.childClassName}`)
      if (childElement) {
        this.fixed_element = childElement
      } else {
        console.log('class name of child node not exist, please check it')
        this.fixed_element = element
      }
    } else {
      this.fixed_element = element
    }
    this.fixed_element_height = this.fixed_element.offsetHeight
    if (this.widthSelf) {
      this.fixed_element_width = this.fixed_element.offsetWidth
      window.addEventListener('resize', _.throttle(this.resizeHandler, 50))
    }
  },
  unmounted () {
    window.removeEventListener('scroll', this.handleScroll)
    if (this.widthSelf) {
      window.removeEventListener('resize', this.resizeHandler)
    }
  },
  methods: {
    resizeHandler () {
      this.fixed_element_width = this.fixed_element.offsetWidth
    },
    getElementToPageTop (el) {
      let top = 0
      while (el !== window.document.body && el) {
        top += el.offsetTop
        el = el.offsetParent
      }
      top += this.fixed_element_height
      return top
    },
    addFixedClass (dom, cls) {
      const arrClassName = dom.className.split(' ')
      const index = this.arrIndexOf(arrClassName, cls)
      if (index === -1) {
        dom.className += ' ' + cls
      }
    },
    removeFixedClass (dom, cls) {
      if (dom.className !== '') {
        const arrClassName = dom.className.split(' ')
        const index = this.arrIndexOf(arrClassName, cls)
        if (index !== -1) {
          arrClassName.splice(index, 1)
          dom.className = arrClassName.join(' ')
        }
      }
    },
    arrIndexOf (arr, v) {
      for (let i = 0; i < arr.length; i++) {
        if (arr[i] === v) {
          return i
        }
      }
      return -1
    },
    handleScroll () {
      const scrollTop = window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop
      this.fixed = scrollTop + window.innerHeight < this.getElementToPageTop(this.fixed_element) // 有可能组件增高或降低，需实时读取高度
    }
  }
}
</script>
