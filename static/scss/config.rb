require 'compass/import-once/activate'
require 'sass-globbing'
require 'susy'
require 'breakpoint'
# Require any additional compass plugins here.

# Target dir
css_dir = "../css"

# Source dir
sass_dir = "./"

# Images dir
images_dir = "../img"

# JavaScript dir
#javascripts_dir = "build/js"

# You can select your preferred output style here (can be overridden via the command line):
output_style = :expanded # or :nested or :compact or :compressed

# To enable relative paths to assets via compass helper functions. Uncomment:
relative_assets = true

# Enable Sourcemap
sourcemap = true

# To disable debugging comments that display the original location of your selectors. Uncomment:
line_comments = false

# Turn off cache busting random query strings
asset_cache_buster :none



# If you prefer the indented syntax, you might want to regenerate this
# project again passing --syntax sass, or you can uncomment this:
# preferred_syntax = :sass
# and then run:
# sass-convert -R --from scss --to sass sass scss && rm -rf sass && mv scss sass
