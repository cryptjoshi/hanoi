import React from 'react';

interface PlaceholderContentProps {
  children: React.ReactNode;
  fullWidthOnLargeScreen?: boolean;
}

const PlaceholderContent: React.FC<PlaceholderContentProps> = ({ 
  children, 
  fullWidthOnLargeScreen = false 
}) => {
  const widthClass = fullWidthOnLargeScreen ? 'w-full' : 'w-full md:w-auto';

  return (
    <div className={`${widthClass} p-4 border rounded-lg bg-muted`}>
      {children}
    </div>
  );
};

export default PlaceholderContent;
