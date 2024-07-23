import { useEffect, useState } from "react";

import { Space, Typography, Image, UploadFile, Form, FormInstance } from "antd";
import { RcFile } from "antd/es/upload";
import { UploadChangeParam } from "antd/lib/upload";
import FormData from "form-data";
import { useTranslation } from "next-i18next";

import EditPencilButton from "../../polaris/button/editPencil";
import DropZoneCard from "./DropZoneCard";
import ExistingContentFilesModal from "./ExistingContentFilesModal";
import { MediaCardContainer } from "./styled";
import FormCard from "@/components/common/form/formCard";
import { PlusIcon } from "@/components/routes/member/announcement/details/styled";
import useSubmit from "@/hooks/common/useSubmit";
import { useGetFiles } from "@/hooks/services/contents/files/useGetFiles";
import useUploadContentFiles from "@/hooks/services/contents/files/useUploadContentFiles";
import { IContentFiles, IUploadFilesData } from "@/models/contents/files";
import { toArrayFormData } from "@/utils/modelToFormData";

interface Props {
  singleUpload?: boolean;
  currentImage?: string;
  form?: FormInstance;
  imageLink?: string;
  detailImage: string;
}

function MediaCard({
  singleUpload = false,
  currentImage,
  form,
  imageLink,
  detailImage,
}: Props) {
  const { t } = useTranslation("common");
  const [visible, setVisible] = useState(false);
  const [selectedKey, setSelectedKey] = useState<string[]>([]);
  const [selectedFiles, setSelectedFiles] = useState<IContentFiles[]>([]);
  const [recentUploadedFiles, setRecentUploadedFiles] = useState<string[]>([]);
  const image = form?.getFieldValue("logo");

  console.log(image);

  const initQuery = {
    folder: "test",
    page: "",
    per_page: "",
  };

  const { data, onReset: onResetGetFiles } = useGetFiles({
    initQuery,
    shouldCall: !!recentUploadedFiles,
  });

  const { doRequest: uploadContentFiles } = useUploadContentFiles();
  const { onFinish: onFinishupload, isLoading } = useSubmit({
    trigger: uploadContentFiles,
    messageSuccess: t("files_uploaded"),
    onSuccess: async () => {
      onResetGetFiles(initQuery);
    },
  });

  const handleFinish = (
    values: IUploadFilesData | UploadChangeParam<UploadFile>
  ) => {
    const files = values?.fileList
      ?.filter((file: UploadFile) => file?.originFileObj !== undefined)
      ?.map((file: UploadFile) => file?.originFileObj);
    if (files && files.length > 0) {
      const formData = toArrayFormData(files as RcFile[]);

      const fileNames = files.map((file: RcFile) => file.name);
      setRecentUploadedFiles(fileNames);

      formData.append("folder", "/test");
      onFinishupload(formData as FormData);
      console.log("finish");

      // form.resetFields();
    }
  };

  useEffect(() => {
    if (data) {
      const selectedCheckFiles = data?.result.filter((item) =>
        selectedKey.includes(item.object_name)
      );
      setSelectedFiles(selectedCheckFiles as IContentFiles[]);

      if (recentUploadedFiles.length > 0) {
        const recentUploadedFilesCheck = data?.result.filter((item) =>
          recentUploadedFiles.includes(item.object_name)
        );

        setSelectedFiles((prev) => [
          ...prev,
          ...(recentUploadedFilesCheck as IContentFiles[]),
        ]);
        setSelectedKey((prev) => [
          ...prev,
          ...(recentUploadedFilesCheck as IContentFiles[]).map(
            (item) => item.object_name
          ),
        ]);
        if (singleUpload) {
          setSelectedFiles([recentUploadedFilesCheck[0] as IContentFiles]);
          setSelectedKey([recentUploadedFiles[0]]);
          form?.setFieldValue("logo", recentUploadedFiles[0]);
        }

        setRecentUploadedFiles([]);
      }
    }
    const existingImage = data?.result.filter(
      (item) => detailImage === item.path
    );

    console.log(
      existingImage,
      currentImage,
      imageLink,
      existingImage?.[0]?.path,
      "cond"
    );
    console.log(form?.getFieldValue("logo"), "logo");
    if (
      existingImage &&
      currentImage &&
      detailImage === form?.getFieldValue("logo")
    ) {
      console.log("run existing");
      setSelectedFiles(existingImage);
      setSelectedKey([currentImage]);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [data, image, form, currentImage]);

  return (
    <>
      <FormCard>
        <Space
          direction="vertical"
          style={{ width: "100%", minHeight: "var(--bo-space-3200)" }}
        >
          <Space style={{ width: "100%", justifyContent: "space-between" }}>
            <Typography.Text strong>{t("image")}</Typography.Text>
            {singleUpload && selectedFiles.length !== 0 && (
              <EditPencilButton onClick={() => setVisible(true)} />
            )}
          </Space>
          <DropZoneCard
            singleUpload={singleUpload}
            hasFile={selectedFiles.length > 0}
            onUpload={(value: UploadChangeParam<UploadFile>) =>
              handleFinish(value)
            }
            onBoxClick={() => setVisible(true)}
            selectedFiles={
              <MediaCardContainer
                $single={singleUpload}
                onClick={(e) => e.stopPropagation()}
              >
                {selectedFiles.map((item) => (
                  <FormCard key={item.object_name}>
                    {/* eslint-disable-next-line jsx-a11y/alt-text */}
                    <Image width={50} height={50} src={item.path} />
                  </FormCard>
                ))}
                <FormCard
                  cardProps={{
                    style: {
                      display: singleUpload ? "none" : "",
                      width: "100%",
                      height: "100%",
                      backgroundColor: "var(--bo-color-bg)",
                      cursor: "pointer",
                    },
                    onClick: () => {
                      setVisible(true);
                    },
                  }}
                >
                  <Space
                    style={{
                      width: "100%",
                      height: "100%",
                      justifyContent: "center",
                    }}
                  >
                    <PlusIcon />
                  </Space>
                </FormCard>
              </MediaCardContainer>
            }
          />
        </Space>
      </FormCard>
      <ExistingContentFilesModal
        single={singleUpload}
        isLoading={isLoading}
        recentUploadedFiles={recentUploadedFiles}
        onUpload={(value: UploadChangeParam<UploadFile>) => handleFinish(value)}
        setOnSelectedFiles={(key) => {
          setSelectedKey(key);
          console.log(key, "key");
          form?.setFieldValue(
            "logo",
            `${process.env.NEXT_PUBLIC_FILE_CDN_BASE_URL}/test/${key}`
          );
          console.log(form?.getFieldValue("logo"), " on change log");
        }}
        onCancel={() => {
          setVisible(false);
          onResetGetFiles(initQuery);
        }}
        visible={visible}
      />
    </>
  );
}

export default MediaCard;
