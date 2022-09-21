cmake_minimum_required(VERSION 3.0)

project(ethz_apriltag2)

#find_package(catkin REQUIRED COMPONENTS cmake_modules)
#include_directories(${catkin_INCLUDE_DIRS})

#catkin_package(
#    DEPENDS eigen opencv
#
#    LIBRARIES ${PROJECT_NAME}
#)

set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} ${PROJECT_SOURCE_DIR}/cmake/)

find_package(Eigen3 PATHS "C:/tools/Eigen3/share/eigen3/cmake")
if (EIGEN3_FOUND)
    INCLUDE_DIRECTORIES(${EIGEN3_INCLUDE_DIRS})
    ADD_DEFINITIONS(${EIGEN3_DEFINITIONS} -D_USE_EIGEN)
    SET(_USE_EIGEN TRUE)
    #    message('${EIGEN3_INCLUDE_DIRS}')
else ()
    message('Eigen not found')
endif ()
set(OpenCV_DIR "C:/opencv/build/install/x64/mingw/lib")
find_package(OpenCV PATHS C:/opencv/build/install/  REQUIRED)
if (OpenCV_FOUND)
    INCLUDE_DIRECTORIES(${OpenCV_INCLUDE_DIRS})
    ADD_DEFINITIONS(${OpenCV_DEFINITIONS})
    SET(_USE_OPENCV TRUE)
    MESSAGE(STATUS "OpenCV ${OpenCV_VERSION} found (include: ${OpenCV_INCLUDE_DIRS})")
    LINK_DIRECTORIES(
        C:/opencv/build/install/x64/mingw/lib
    )
else ()
    message("OpenCV not found.")
endif ()


add_definitions(-fPIC -O3)
include_directories(include ${Eigen_INCLUDE_DIRS})

#library
file(GLOB SOURCE_FILES "src/*.cc")
add_library(${PROJECT_NAME} SHARED ${SOURCE_FILES})
target_link_libraries(${PROJECT_NAME} ${OpenCV_LIBS})

#demo
if (NOT APPLE)
    add_executable(apriltags_demo src/example/apriltags_demo.cpp src/example/Serial.cpp)
    target_link_libraries(apriltags_demo ${PROJECT_NAME} v4l2)
endif ()


