# cmakefile for check testing
# followed from http://isti.bitbucket.org/2012/06/01/cmake-check.html
#               http://johnlamp.net/cmake-tutorial-4-libraries-and-subdirectories.html
# FindGTest: https://github.com/Kitware/CMake/blob/master/Modules/FindGTest.cmake
enable_testing()
find_package(GTest REQUIRED)
set(PROJECT_TEST_NAME "${PROJECT_NAME}_test")
include_directories(BEFORE ${GTEST_INCLUDE_DIRS})
file(GLOB TEST_SRC_FILES ${PROJECT_SOURCE_DIR}/tests/*.cc)
add_executable(option_tests ${TEST_SRC_FILES})
target_link_libraries(option_tests
    ${PROJECT_NAME_STR}
    ${GTEST_BOTH_LIBRARIES}
    ${Boost_LIBS}
    pthread
    )
add_test(AllTests option_tests)
