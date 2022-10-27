from tkinter import image_types
import cv2
import os

# img_path = "blur_test.jpg"

data_dir = "C:/Users/liudx/Desktop/record_C08"

image_paths = [os.path.join(data_dir, f)
               for f in os.listdir(data_dir) if f.endswith(".png")]
cv2.namedWindow("blurDetect", cv2.WINDOW_NORMAL)
for img_path in image_paths:
    gray = cv2.imread(img_path, cv2.IMREAD_GRAYSCALE)
    before_blur = cv2.Laplacian(gray, cv2.CV_64F).var()  # get blur value

    img = cv2.imread(img_path)  # load image
    # for i in range(20):  # blur the picture
    #     img = cv2.blur(img, (3, 3))
    # gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)  # bgr to gray
    # after_blur = cv2.Laplacian(gray, cv2.CV_64F).var()  # get blur value

    # put text to blur image
    text1 = f'image shape: {gray.shape}'
    text2 = f'before blur: {before_blur}'
    # text3 = f'after blur: {after_blur}'
    if True:
        cv2.putText(img, text1, (10, 30),
                    cv2.FONT_HERSHEY_SIMPLEX, 0.8, (0, 255, 0), 2)
        cv2.putText(img, text2, (10, 60),
                    cv2.FONT_HERSHEY_SIMPLEX, 0.8, (0, 255, 0), 2)
        # cv2.putText(img, text3, (10, 90),
        #             cv2.FONT_HERSHEY_SIMPLEX, 0.8, (0, 255, 0), 2)

        # show image
        cv2.imshow('blurDetect', img)
        if before_blur < 490:
            cv2.waitKey(0)
        cv2.waitKey(2)
    else:
        cv2.imshow('blurDetect', img)
        cv2.waitKey(1)
cv2.destroyAllWindows()
