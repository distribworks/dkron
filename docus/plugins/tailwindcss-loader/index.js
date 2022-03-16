module.exports = function (context, options) {
  return {
    name: 'postcss-tailwindcss-loader',
    configurePostCss(postcssOptions) {
      postcssOptions.plugins.push(
        require('postcss-import'),
        require('tailwindcss'),
      )
      return postcssOptions
    },
  }
}