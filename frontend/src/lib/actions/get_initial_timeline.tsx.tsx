"use client";
import useSWR from "swr";
import { GetCollectionOfPostsBySpecificUserAndUsersTheyFollowResponse } from "./get-collection-of-posts-by-specific-user-and-users-they-follow";

type ApiResponse = GetCollectionOfPostsBySpecificUserAndUsersTheyFollowResponse;

type ApiError = {
  message: string;
  status: number;
};

export const getInitialTimeline = () => {
  const sample_user_id = "012";

  const fetcher = async (url: string) => {
    const response = await fetch(url, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error: ApiError = {
        message: `An error occurred: ${response.statusText}`,
        status: response.status,
      };
      throw error;
    }

    const data = await response.json();
    if (!data) {
      return [];
    }
    return data;
  };

  const { data, error } = useSWR<ApiResponse, ApiError>(
    `http://localhost:80/api/${sample_user_id}/polling?event_type=TimelineAccessed`,
    fetcher,
    {
      revalidateOnFocus: false,
      revalidateOnReconnect: false,
    }
  );

  return { data, error };
};
