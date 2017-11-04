var gulp = require('gulp')
var riot = require('gulp-riot')
var browserify = require('gulp-browserify')
var source = require('vinyl-source-stream')
var path = require('path')

const paths = {
    scripts:['tags/**/*.tag']
}

gulp.task('scripts', function() {
    return gulp.src(paths.scripts)
        .pipe(riot())
        .pipe(gulp.dest('tags'))
})

gulp.task('browserify', function() {
    gulp.src('tags/main.js')
        .pipe(browserify())
        .pipe(gulp.dest('static/js'))
})

gulp.task('watch', function() {
    gulp.watch(paths.scripts, ['scripts', 'browserify'])
})

gulp.task('default', ['watch', 'scripts', 'browserify'])